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

export class Version {
  majorDigit;
  minorDigit;
  patchDigit;
  suffix;

  constructor(version) {
    const [majorDigit, minorDigit, patchDigit, suffix] = [
      ...version.split(".").map(Number),
    ];
    Object.assign(this, { majorDigit, minorDigit, patchDigit, suffix });
  }

  branch() {
    return `${this.majorDigit}.${this.minorDigit}.x`;
  }

  minor() {
    return `${this.majorDigit}.${this.minorDigit}`;
  }

  isNotPatch() {
    return this.patchDigit === 0;
  }

  isPreRelease() {
    return !!this.suffix;
  }

  next() {
    return this.isNotPatch() ? this.nextMinor() : this.nextPatch();
  }

  nextMinor() {
    return new Version(
      `${this.majorDigit}.${this.minorDigit + 1}.${this.patchDigit}`,
    );
  }

  nextPatch() {
    return new Version(
      `${this.majorDigit}.${this.minorDigit}.${this.patchDigit + 1}`,
    );
  }

  toString() {
    return `${this.majorDigit}.${this.minorDigit}.${this.patchDigit}`;
  }
}
