/**
 * Logging class that wraps logs
 */
export default class Logger {
  static progressColor = "\x1b[36m%s\x1b[0m"; // Cyan
  static successColor = "\x1b[32m%s\x1b[0m"; // Light green
  static errorColor = "\x1b[31m%s\x1b[0m"; // Red

  /**
   * Log an info message
   * @param message Info message
   * @example
   * // Logs 'Info message'
   * logger.info('Info message');
   */
  static info(message: string) {
    console.log(message);
  }

  /**
   * Log a progress message
   * @param message Progress message
   * @example
   * // Logs 'Progress message...' in progress color
   * logger.progress('Progress message');
   */
  static progress(message: string) {
    console.log(Logger.progressColor, `${message}...`);
  }

  /**
   * Log a success message
   * @param message Success message
   * @example
   * // Logs 'Success message' in success color
   * success('Success message');
   */
  static success(message: string) {
    console.log(Logger.successColor, message);
  }

  /**
   * Log an error message
   * @param message Error message
   * @example
   * // Logs 'Error message' in error color
   * error('Error message');
   */
  static error(message: string) {
    console.error(Logger.errorColor, message);
  }
}
