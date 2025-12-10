#!/usr/bin/env zx
/**
 * Copyright (C) 2015 The Gravitee team (http://gravitee.io)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import path from 'node:path';
import { readFile } from 'node:fs/promises';
import { isDeepStrictEqual } from 'node:util';
import { parse as parseYaml } from 'yaml';

const firstPath = argv.first ?? argv.f;
const secondPath = argv.second ?? argv.s;
const ignorePaths = normalizeIgnorePaths(argv.ignore ?? argv.i);

if (!firstPath || !secondPath) {
  console.error('Usage: compareYamlFiles.mjs --first <file> --second <file> [--ignore <path>]...');
  process.exit(1);
}

const firstAbsolute = path.resolve(process.cwd(), firstPath);
const secondAbsolute = path.resolve(process.cwd(), secondPath);

const firstObject = await loadYaml(firstAbsolute, 'first document');
const secondObject = await loadYaml(secondAbsolute, 'second document');

const shouldIgnore = createIgnoreChecker(ignorePaths);
const differences = findDifferences(firstObject, secondObject, '$', shouldIgnore);

if (!differences.length) {
  console.log('✅ YAML documents are identical.');
  process.exit(0);
}

console.error('❌ YAML documents differ:');
for (const message of differences) {
  console.error(` - ${message}`);
}
console.error(`Total mismatches: ${differences.length}.`);
process.exit(1);

async function loadYaml(filePath, label) {
  try {
    const raw = await readFile(filePath, 'utf8');
    if (!raw.trim()) {
      throw new Error('file is empty');
    }
    return parseYaml(raw);
  } catch (error) {
    console.error(`Failed to read ${label} from ${filePath}: ${error.message}`);
    process.exit(1);
  }
}

function findDifferences(firstValue, secondValue, currentPath = '$', shouldIgnore = () => false) {
  if (shouldIgnore(currentPath)) {
    return [];
  }

  if (isDeepStrictEqual(firstValue, secondValue)) {
    return [];
  }

  if (firstValue === undefined) {
    return [`${currentPath}: missing in first document but present in second (${display(secondValue)})`];
  }

  if (secondValue === undefined) {
    return [`${currentPath}: present in first document (${display(firstValue)}) but missing in second`];
  }

  const firstIsArray = Array.isArray(firstValue);
  const secondIsArray = Array.isArray(secondValue);
  if (firstIsArray && secondIsArray) {
    const issues = [];
    if (firstValue.length !== secondValue.length) {
      issues.push(`${currentPath}: array length differs (first=${firstValue.length}, second=${secondValue.length})`);
    }

    const max = Math.max(firstValue.length, secondValue.length);
    for (let index = 0; index < max; index += 1) {
      issues.push(...findDifferences(firstValue[index], secondValue[index], `${currentPath}[${index}]`, shouldIgnore));
    }
    return issues;
  }

  const firstIsObject = isPlainObject(firstValue);
  const secondIsObject = isPlainObject(secondValue);
  if (firstIsObject && secondIsObject) {
    const keys = new Set([...Object.keys(firstValue), ...Object.keys(secondValue)]);
    const messages = [];
    for (const key of [...keys].sort()) {
      const nextPath = buildChildPath(currentPath, key);
      messages.push(...findDifferences(firstValue[key], secondValue[key], nextPath, shouldIgnore));
    }
    return messages;
  }

  if (firstIsArray !== secondIsArray || firstIsObject !== secondIsObject) {
    return [`${currentPath}: type mismatch (first=${describe(firstValue)}, second=${describe(secondValue)})`];
  }

  return [`${currentPath}: value mismatch (first=${display(firstValue)}, second=${display(secondValue)})`];
}

function isPlainObject(value) {
  return typeof value === 'object' && value !== null && !Array.isArray(value);
}

function buildChildPath(parent, key) {
  if (/^[A-Za-z0-9_]+$/.test(key)) {
    return parent === '$' ? `${parent}.${key}` : `${parent}.${key}`;
  }
  return `${parent}[${JSON.stringify(key)}]`;
}

function normalizeIgnorePaths(rawValue) {
  return toArray(rawValue)
    .map((entry) => String(entry).trim())
    .filter(Boolean)
    .map((entry) => (entry.startsWith('$') ? entry : `$.${entry}`));
}

function toArray(value) {
  if (value === undefined || value === null) {
    return [];
  }
  return Array.isArray(value) ? value : [value];
}

function createIgnoreChecker(paths) {
  if (!paths.length) {
    return () => false;
  }
  const matchers = paths.map((path) => (path === '$' ? '$' : path.replace(/\.$/, '')));
  return (currentPath) =>
    matchers.some((ignored) => currentPath === ignored || currentPath.startsWith(`${ignored}.`) || currentPath.startsWith(`${ignored}[`));
}

function describe(value) {
  if (Array.isArray(value)) {
    return 'array';
  }
  if (isPlainObject(value)) {
    return 'object';
  }
  return typeof value;
}

function display(value) {
  if (typeof value === 'string') {
    return JSON.stringify(value);
  }
  if (value === null) {
    return 'null';
  }
  if (value === undefined) {
    return 'undefined';
  }
  try {
    return JSON.stringify(value);
  } catch {
    return String(value);
  }
}
