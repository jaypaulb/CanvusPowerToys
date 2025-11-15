/**
 * Error Handler - Centralized error handling for WebUI
 * Replaces console.error with proper error handling
 */

class ErrorHandler {
  constructor() {
    this.isDevelopment = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1';
  }

  /**
   * Log error with appropriate handling
   * @param {string} message - Error message
   * @param {Error} error - Error object
   * @param {string} context - Context where error occurred
   */
  logError(message, error, context = '') {
    const errorMessage = context ? `${context}: ${message}` : message;
    
    // In development, log to console for debugging
    if (this.isDevelopment) {
      console.error(errorMessage, error);
    }
    
    // In production, could send to error tracking service
    // For now, silently handle (or show user-friendly message)
    this.showUserError(message);
  }

  /**
   * Show user-friendly error message
   * @param {string} message - User-friendly error message
   */
  showUserError(message) {
    // Could show toast notification or update UI
    // For now, errors are handled silently in production
    // Individual pages can override this behavior
  }
}

// Export singleton instance
const errorHandler = new ErrorHandler();

// Export for use in other modules
if (typeof module !== 'undefined' && module.exports) {
  module.exports = errorHandler;
}

