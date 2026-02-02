#!/usr/bin/env node

const childProcess = require('child_process');

const BINARY_MAP = {
  darwin_x64: {name: 'mcp-prometheus-darwin-amd64', suffix: ''},
  darwin_arm64: {name: 'mcp-prometheus-darwin-arm64', suffix: ''},
  linux_x64: {name: 'mcp-prometheus-linux-amd64', suffix: ''},
  linux_arm64: {name: 'mcp-prometheus-linux-arm64', suffix: ''},
  win32_x64: {name: 'mcp-prometheus-windows-amd64', suffix: '.exe'},
  win32_arm64: {name: 'mcp-prometheus-windows-arm64', suffix: '.exe'},
};

const resolveBinaryPath = () => {
  try {
    const binary = BINARY_MAP[`${process.platform}_${process.arch}`];
    return require.resolve(`${binary.name}/bin/${binary.name}${binary.suffix}`);
  } catch (e) {
    throw new Error(`Could not resolve binary path for platform/arch: ${process.platform}/${process.arch}`);
  }
};

const child = childProcess.spawn(resolveBinaryPath(), process.argv.slice(2), {
  stdio: 'inherit',
});

const handleSignal = (signal) => {
  if (child && !child.killed) {
    child.kill(signal);
  }
};

['SIGTERM', 'SIGINT', 'SIGHUP'].forEach((signal) => {
  process.on(signal, handleSignal);
});

child.on('close', (code, signal) => {
  if (signal) {
    process.exit(128 + (signal === 'SIGTERM' ? 15 : signal === 'SIGINT' ? 2 : 1));
  } else {
    process.exit(code || 0);
  }
});
