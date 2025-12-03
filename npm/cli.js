#!/usr/bin/env node

const { spawn, execSync } = require('child_process');
const fs = require('fs');
const path = require('path');
const https = require('https');

const GITHUB_REPO = 'jeanlopezxyz/mcp-prometheus';
const JAR_NAME = 'mcp-prometheus.jar';
const CACHE_DIR = path.join(require('os').homedir(), '.cache', 'mcp-prometheus');

function parseArgs() {
  const args = process.argv.slice(2);
  const options = { port: null, help: false, version: false, extraArgs: [] };

  for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    if (arg === '--port' && args[i + 1]) {
      options.port = args[++i];
    } else if (arg.startsWith('--port=')) {
      options.port = arg.split('=')[1];
    } else if (arg === '--help' || arg === '-h') {
      options.help = true;
    } else if (arg === '--version' || arg === '-v') {
      options.version = true;
    } else {
      options.extraArgs.push(arg);
    }
  }
  return options;
}

function showHelp() {
  console.log(`
mcp-prometheus - MCP Server for Prometheus

USAGE:
  npx mcp-prometheus [OPTIONS]

OPTIONS:
  --port <PORT>    Start in SSE mode on specified port (default: stdio mode)
  --help, -h       Show this help message
  --version, -v    Show version

ENVIRONMENT:
  PROMETHEUS_URL   Prometheus API URL (default: http://localhost:9090)

EXAMPLES:
  # stdio mode (for Claude Code, Claude Desktop)
  PROMETHEUS_URL="http://prometheus:9090" npx mcp-prometheus

  # SSE mode on port 9081
  PROMETHEUS_URL="http://prometheus:9090" npx mcp-prometheus --port 9081

CLAUDE CODE CONFIGURATION:
  Add to ~/.claude/settings.json:

  {
    "mcpServers": {
      "prometheus": {
        "command": "npx",
        "args": ["-y", "mcp-prometheus@latest"],
        "env": {
          "PROMETHEUS_URL": "http://localhost:9090"
        }
      }
    }
  }
`);
}

function checkJava() {
  try {
    const result = execSync('java -version 2>&1', { encoding: 'utf8' });
    const match = result.match(/version "(\d+)/);
    if (match && parseInt(match[1]) < 21) {
      console.error('Error: Java 21+ is required. Found:', result.split('\n')[0]);
      process.exit(1);
    }
    return true;
  } catch {
    console.error('Error: Java 21+ is required but not found.');
    console.error('Install Java from: https://adoptium.net/');
    process.exit(1);
  }
}

async function getLatestRelease() {
  return new Promise((resolve, reject) => {
    const options = {
      hostname: 'api.github.com',
      path: `/repos/${GITHUB_REPO}/releases/latest`,
      headers: { 'User-Agent': 'mcp-prometheus' }
    };

    https.get(options, (res) => {
      let data = '';
      res.on('data', chunk => data += chunk);
      res.on('end', () => {
        try {
          resolve(JSON.parse(data));
        } catch (e) {
          reject(new Error('Failed to parse release info'));
        }
      });
    }).on('error', reject);
  });
}

function downloadFile(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);

    const request = (url) => {
      https.get(url, { headers: { 'User-Agent': 'mcp-prometheus' } }, (res) => {
        if (res.statusCode === 302 || res.statusCode === 301) {
          request(res.headers.location);
          return;
        }
        if (res.statusCode !== 200) {
          reject(new Error(`Download failed: ${res.statusCode}`));
          return;
        }
        res.pipe(file);
        file.on('finish', () => { file.close(); resolve(); });
      }).on('error', (err) => { fs.unlink(dest, () => {}); reject(err); });
    };

    request(url);
  });
}

async function getJar(verbose = true) {
  const jarPath = path.join(CACHE_DIR, JAR_NAME);
  const versionFile = path.join(CACHE_DIR, 'version');

  if (!fs.existsSync(CACHE_DIR)) {
    fs.mkdirSync(CACHE_DIR, { recursive: true });
  }

  try {
    const release = await getLatestRelease();
    const latestVersion = release.tag_name;

    let currentVersion = null;
    if (fs.existsSync(versionFile)) {
      currentVersion = fs.readFileSync(versionFile, 'utf8').trim();
    }

    if (currentVersion === latestVersion && fs.existsSync(jarPath)) {
      return jarPath;
    }

    const jarAsset = release.assets.find(a => a.name.endsWith('.jar'));
    if (!jarAsset) {
      throw new Error('No JAR found in release');
    }

    if (verbose) console.error(`Downloading mcp-prometheus ${latestVersion}...`);
    await downloadFile(jarAsset.browser_download_url, jarPath);
    fs.writeFileSync(versionFile, latestVersion);
    if (verbose) console.error('Download complete.');

    return jarPath;
  } catch (error) {
    if (fs.existsSync(jarPath)) {
      if (verbose) console.error('Warning: Could not check for updates, using cached version.');
      return jarPath;
    }
    throw error;
  }
}

async function main() {
  const options = parseArgs();

  if (options.help) { showHelp(); process.exit(0); }
  if (options.version) { console.log(require('./package.json').version); process.exit(0); }

  if (!process.env.PROMETHEUS_URL) {
    process.env.PROMETHEUS_URL = 'http://localhost:9090';
  }

  checkJava();

  try {
    const verbose = !!options.port;
    const jarPath = await getJar(verbose);

    const javaArgs = [];

    if (options.port) {
      javaArgs.push(`-Dquarkus.http.port=${options.port}`);
      javaArgs.push('-Dquarkus.http.host=0.0.0.0');
      console.error(`Starting MCP server in SSE mode on port ${options.port}...`);
      console.error(`SSE endpoint: http://localhost:${options.port}/mcp/sse`);
    } else {
      javaArgs.push('-Dquarkus.http.host-enabled=false');
      javaArgs.push('-Dquarkus.banner.enabled=false');
      javaArgs.push('-Dquarkus.log.level=WARN');
      javaArgs.push('-Dquarkus.mcp.server.traffic-logging.enabled=false');
    }

    javaArgs.push('-jar', jarPath);
    javaArgs.push(...options.extraArgs);

    const java = spawn('java', javaArgs, { stdio: 'inherit', env: process.env });

    java.on('error', (err) => { console.error('Failed to start Java:', err.message); process.exit(1); });
    java.on('exit', (code) => { process.exit(code || 0); });

    const handleSignal = (signal) => { if (java && !java.killed) java.kill(signal); };
    process.on('SIGINT', () => handleSignal('SIGINT'));
    process.on('SIGTERM', () => handleSignal('SIGTERM'));

  } catch (error) {
    console.error('Error:', error.message);
    process.exit(1);
  }
}

main();
