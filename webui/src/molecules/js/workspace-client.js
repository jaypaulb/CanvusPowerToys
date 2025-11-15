/**
 * Workspace Client Molecule - SSE Client Utility
 * Handles connection to /api/subscribe-workspace endpoint
 */

class WorkspaceClient {
  constructor() {
    this.eventSource = null;
    this.listeners = new Map();
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectDelay = 3000;
    this.isConnected = false;
  }

  /**
   * Connect to workspace subscription endpoint
   * @param {string} baseURL - Base URL for the API
   */
  connect(baseURL) {
    if (this.eventSource) {
      this.disconnect();
    }

    const url = `${baseURL}/api/subscribe-workspace`;
    this.eventSource = new EventSource(url);

    this.eventSource.onopen = () => {
      this.isConnected = true;
      this.reconnectAttempts = 0;
      this.emit('connected');
    };

    this.eventSource.addEventListener('canvas_update', (event) => {
      try {
        const data = JSON.parse(event.data);
        this.emit('canvas_update', data);
      } catch (error) {
        // Only log in development
        if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
          console.error('Error parsing canvas update:', error);
        }
      }
    });

    this.eventSource.onerror = (error) => {
      this.isConnected = false;
      this.emit('error', error);

      if (this.eventSource.readyState === EventSource.CLOSED) {
        this.handleReconnect(baseURL);
      }
    };
  }

  /**
   * Disconnect from workspace subscription
   */
  disconnect() {
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = null;
      this.isConnected = false;
      this.emit('disconnected');
    }
  }

  /**
   * Handle reconnection logic
   * @param {string} baseURL - Base URL for the API
   */
  handleReconnect(baseURL) {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      this.emit('max_reconnect_attempts');
      return;
    }

    this.reconnectAttempts++;
    this.emit('reconnecting', this.reconnectAttempts);

    setTimeout(() => {
      this.connect(baseURL);
    }, this.reconnectDelay);
  }

  /**
   * Subscribe to an event
   * @param {string} event - Event name
   * @param {Function} callback - Callback function
   */
  on(event, callback) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, []);
    }
    this.listeners.get(event).push(callback);
  }

  /**
   * Unsubscribe from an event
   * @param {string} event - Event name
   * @param {Function} callback - Callback function to remove
   */
  off(event, callback) {
    if (this.listeners.has(event)) {
      const callbacks = this.listeners.get(event);
      const index = callbacks.indexOf(callback);
      if (index > -1) {
        callbacks.splice(index, 1);
      }
    }
  }

  /**
   * Emit an event to all listeners
   * @param {string} event - Event name
   * @param {*} data - Event data
   */
  emit(event, data) {
    if (this.listeners.has(event)) {
      this.listeners.get(event).forEach(callback => {
        try {
          callback(data);
        } catch (error) {
          // Only log in development
          if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
            console.error(`Error in event listener for ${event}:`, error);
          }
        }
      });
    }
  }
}

// Export for use in other modules
if (typeof module !== 'undefined' && module.exports) {
  module.exports = WorkspaceClient;
}

