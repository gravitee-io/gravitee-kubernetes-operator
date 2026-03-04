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

import { describe, it, expect, vi, afterEach } from "vitest";
import { readFile } from "node:fs/promises";
import { loadMatchFile } from "../../src/cmd/assert-api.js";

vi.mock("node:fs/promises", () => ({
  readFile: vi.fn(),
}));

const mockedReadFile = vi.mocked(readFile);

afterEach(() => {
  vi.restoreAllMocks();
});

describe("loadMatchFile", () => {
  it("parses a valid YAML file", async () => {
    mockedReadFile.mockResolvedValue("state: STARTED\nname: My API\n");
    const result = await loadMatchFile("/path/to/expect.yaml");
    expect(result).toEqual({ state: "STARTED", name: "My API" });
    expect(mockedReadFile).toHaveBeenCalledWith("/path/to/expect.yaml", "utf-8");
  });

  it("parses nested YAML structures", async () => {
    mockedReadFile.mockResolvedValue(
      "listeners:\n  - type: HTTP\n    paths:\n      - path: /petstore\ncategories:\n  - finance\n",
    );
    const result = await loadMatchFile("/path/to/expect.yaml");
    expect(result).toEqual({
      listeners: [{ type: "HTTP", paths: [{ path: "/petstore" }] }],
      categories: ["finance"],
    });
  });

  it("throws when the file does not exist", async () => {
    mockedReadFile.mockRejectedValue(new Error("ENOENT: no such file or directory"));
    await expect(loadMatchFile("/missing.yaml")).rejects.toThrow(
      'assert-api: cannot read --match-file "/missing.yaml"',
    );
  });

  it("throws when the file contains invalid YAML", async () => {
    mockedReadFile.mockResolvedValue("key: [unterminated");
    await expect(loadMatchFile("/bad.yaml")).rejects.toThrow(
      'assert-api: --match-file "/bad.yaml" is not valid YAML',
    );
  });

  it("throws when the file contains a YAML scalar instead of a mapping", async () => {
    mockedReadFile.mockResolvedValue("just a string");
    await expect(loadMatchFile("/scalar.yaml")).rejects.toThrow(
      'assert-api: --match-file "/scalar.yaml" must contain a YAML mapping (object)',
    );
  });

  it("throws when the file contains a YAML array instead of a mapping", async () => {
    mockedReadFile.mockResolvedValue("- item1\n- item2\n");
    await expect(loadMatchFile("/array.yaml")).rejects.toThrow(
      'assert-api: --match-file "/array.yaml" must contain a YAML mapping (object), got array',
    );
  });

  it("throws when the file is empty (parses to null)", async () => {
    mockedReadFile.mockResolvedValue("");
    await expect(loadMatchFile("/empty.yaml")).rejects.toThrow(
      'assert-api: --match-file "/empty.yaml" must contain a YAML mapping (object)',
    );
  });
});
