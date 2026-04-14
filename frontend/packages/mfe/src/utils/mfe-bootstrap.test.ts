import { describe, test, expect } from 'vitest';
import { parsePath } from './mfe-bootstrap.ts';

describe('mfe-bootstrap.ts', () => {
  describe('parsePath', () => {
    test('should return empty string when path is empty', () => {
      // Arrange
      const path = '';
      const basePath = '/base';

      // Act
      const result = parsePath(path, basePath);

      // Assert
      expect(result).toBe('');
    });

    test('should return path when basePath is empty', () => {
      // Arrange
      const path = '/some/path';
      const basePath = '';

      // Act
      const result = parsePath(path, basePath);

      // Assert
      expect(result).toBe('/some/path');
    });

    test('should return original path when path does not start with basePath', () => {
      // Arrange
      const path = '/other/path';
      const basePath = '/base';

      // Act
      const result = parsePath(path, basePath);

      // Assert
      expect(result).toBe('/other/path');
    });

    test('should remove basePath from path when path starts with basePath', () => {
      // Arrange
      const path = '/base/some/path';
      const basePath = '/base';

      // Act
      const result = parsePath(path, basePath);

      // Assert
      expect(result).toBe('/some/path');
    });

    test('should return empty string when path equals basePath', () => {
      // Arrange
      const path = '/base';
      const basePath = '/base';

      // Act
      const result = parsePath(path, basePath);

      // Assert
      expect(result).toBe('');
    });
  });
});
