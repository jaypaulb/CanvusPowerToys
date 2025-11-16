/**
 * RCU Admin Page JavaScript
 * Handles team target creation, test team notes, and user management
 */

document.addEventListener('DOMContentLoaded', () => {
  initTeamButtonStyles();
  initCreateTargets();
  initTestTeamButtons();
  initUserManagement();
});

/**
 * Initialize team button styles (ROYGBIV colors) - FIXED SIZE
 */
function initTeamButtonStyles() {
  const teamButtons = document.querySelectorAll('[data-test-team]');

  // Team colors (ROYGBIV) - exact RGB values as specified
  const teamColors = {
    1: "rgb(255, 0, 0)",      // Red
    2: "rgb(255, 127, 0)",    // Orange
    3: "rgb(255, 255, 0)",    // Yellow
    4: "rgb(0, 255, 0)",      // Green
    5: "rgb(0, 0, 255)",      // Blue
    6: "rgb(75, 0, 130)",     // Indigo
    7: "rgb(139, 0, 255)"     // Violet
  };

  // Text colors for each team (matching the provided style exactly)
  const teamTextColors = {
    1: "rgb(255, 255, 255)",  // White for Red
    2: "rgb(0, 0, 0)",         // Black for Orange
    3: "rgb(0, 0, 0)",         // Black for Yellow
    4: "rgb(0, 0, 0)",         // Black for Green
    5: "rgb(255, 255, 255)",  // White for Blue
    6: "rgb(255, 255, 255)",  // White for Indigo
    7: "rgb(255, 255, 255)"    // White for Violet
  };

  // Apply colors and FIXED sizing - exactly as specified
  teamButtons.forEach(button => {
    const teamNumber = parseInt(button.getAttribute('data-test-team'));
    if (teamNumber && teamColors[teamNumber]) {
      // Set fixed size - CRITICAL: these must be exactly 120px x 40px
      button.style.width = '120px';
      button.style.height = '40px';
      button.style.flexShrink = '0';
      button.style.minWidth = '120px';
      button.style.maxWidth = '120px';

      // Set colors
      button.style.backgroundColor = teamColors[teamNumber];
      button.style.color = teamTextColors[teamNumber];
    }
  });
}

/**
 * Initialize create targets functionality
 */
function initCreateTargets() {
  const createTargetsBtn = document.getElementById('createTargetsBtn');
  const deleteTargetsBtn = document.getElementById('deleteTargetsBtn');
  const targetsMessage = document.getElementById('targetsMessage');

  if (createTargetsBtn) {
    createTargetsBtn.addEventListener('click', async () => {
      try {
        const response = await fetch('/api/admin/create-targets', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          }
        });

        const data = await response.json();
        if (response.ok && data.success) {
          displayMessage(targetsMessage, data.message || "Target notes created successfully", "success");
        } else {
          displayMessage(targetsMessage, data.error || "Failed to create target notes", "error");
        }
      } catch (error) {
        console.error('Error creating targets:', error);
        displayMessage(targetsMessage, "An error occurred while creating targets", "error");
      }
    });
  }

  if (deleteTargetsBtn) {
    deleteTargetsBtn.addEventListener('click', async () => {
      if (!confirm('Are you sure you want to delete all target notes?')) {
        return;
      }

      try {
        const response = await fetch('/api/admin/delete-targets', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          }
        });

        const data = await response.json();
        if (response.ok && data.success) {
          displayMessage(targetsMessage, data.message || "Target notes deleted successfully", "success");
        } else {
          displayMessage(targetsMessage, data.error || "Failed to delete target notes", "error");
        }
      } catch (error) {
        console.error('Error deleting targets:', error);
        displayMessage(targetsMessage, "An error occurred while deleting targets", "error");
      }
    });
  }
}

/**
 * Initialize test team buttons
 */
function initTestTeamButtons() {
  const testTeamButtons = document.querySelectorAll('[data-test-team]');
  const testTeamMessage = document.getElementById('testTeamMessage');

  testTeamButtons.forEach(button => {
    button.addEventListener('click', async () => {
      const teamNumber = parseInt(button.getAttribute('data-test-team'));

      try {
        // Send a test note from Admin to the team's target
        const response = await fetch('/api/admin/test-team', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            team: teamNumber,
            text: `Test note from Admin to Team ${teamNumber}`
          })
        });

        const data = await response.json();
        if (response.ok && data.success) {
          displayMessage(testTeamMessage, `Test note sent to Team ${teamNumber} successfully`, "success");
        } else {
          displayMessage(testTeamMessage, data.error || `Failed to send test note to Team ${teamNumber}`, "error");
        }
      } catch (error) {
        console.error(`Error sending test note to Team ${teamNumber}:`, error);
        displayMessage(testTeamMessage, `An error occurred while sending test note to Team ${teamNumber}`, "error");
      }
    });
  });
}

/**
 * Initialize user management
 */
function initUserManagement() {
  const listUsersBtn = document.getElementById('listUsersBtn');
  const deleteUsersBtn = document.getElementById('deleteUsersBtn');
  const userListTable = document.getElementById('userListTable');
  const userListMessage = document.getElementById('userListMessage');
  const tbody = userListTable ? userListTable.querySelector('tbody') : null;

  if (listUsersBtn) {
    listUsersBtn.addEventListener('click', async () => {
      try {
        const response = await fetch('/api/admin/list-users', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json'
          }
        });

        const data = await response.json();
        if (response.ok && data.success && data.users) {
          displayUserList(data.users, tbody, userListMessage);
        } else {
          displayMessage(userListMessage, data.error || "Failed to list users", "error");
        }
      } catch (error) {
        console.error('Error listing users:', error);
        displayMessage(userListMessage, "An error occurred while listing users", "error");
      }
    });
  }

  if (deleteUsersBtn) {
    deleteUsersBtn.addEventListener('click', async () => {
      if (!confirm('Are you sure you want to delete all users?')) {
        return;
      }

      try {
        const response = await fetch('/api/admin/delete-users', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({ all: true })
        });

        const data = await response.json();
        if (response.ok && data.success) {
          displayMessage(userListMessage, data.message || "Users deleted successfully", "success");
          if (tbody) tbody.innerHTML = '';
        } else {
          displayMessage(userListMessage, data.error || "Failed to delete users", "error");
        }
      } catch (error) {
        console.error('Error deleting users:', error);
        displayMessage(userListMessage, "An error occurred while deleting users", "error");
      }
    });
  }
}

/**
 * Display user list
 */
function displayUserList(users, tbody, messageEl) {
  if (!tbody) return;

  if (!users || users.length === 0) {
    tbody.innerHTML = '';
    if (messageEl) {
      messageEl.textContent = 'No users to display.';
      messageEl.style.display = 'block';
    }
    return;
  }

  if (messageEl) {
    messageEl.style.display = 'none';
  }

  tbody.innerHTML = users.map(user => `
    <tr>
      <td>Team ${user.team}</td>
      <td>${user.name}</td>
      <td>${user.color ? user.color.toUpperCase() : 'N/A'}</td>
    </tr>
  `).join('');
}

/**
 * Display message
 */
function displayMessage(element, text, type) {
  if (!element) {
    console.log(`[${type}] ${text}`);
    return;
  }

  element.textContent = text;
  element.className = `message ${type} mt-md`;
  element.style.display = 'block';

  // Auto-hide after 5 seconds
  setTimeout(() => {
    element.style.display = 'none';
  }, 5000);
}
