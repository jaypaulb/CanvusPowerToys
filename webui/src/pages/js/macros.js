/**
 * Macros Page JavaScript
 * Handles grouping, pinning, and management of macros
 */

// Tab switching
document.addEventListener('DOMContentLoaded', () => {
  const tabButtons = document.querySelectorAll('.tab-btn');
  const tabContents = document.querySelectorAll('.tab-content');

  tabButtons.forEach(button => {
    button.addEventListener('click', () => {
      const targetTab = button.getAttribute('data-tab');

      // Remove active class from all tabs
      tabButtons.forEach(btn => btn.classList.remove('active'));
      tabContents.forEach(content => content.classList.remove('active'));

      // Add active class to selected tab
      button.classList.add('active');
      document.getElementById(`${targetTab}Tab`).classList.add('active');
    });
  });

  // Load macros data
  loadMacrosData();
});

/**
 * Load macros data from API
 */
async function loadMacrosData() {
  const baseURL = window.location.origin;

  try {
    // Load groups
    const groupsResponse = await fetch(`${baseURL}/api/macros/groups`);
    if (groupsResponse.ok) {
      const groups = await groupsResponse.json();
      renderGroups(groups);
    }

    // Load pinned macros
    const pinnedResponse = await fetch(`${baseURL}/api/macros/pinned`);
    if (pinnedResponse.ok) {
      const pinned = await pinnedResponse.json();
      renderPinnedMacros(pinned);
    }

    // Load all macros for manage tab
    const macrosResponse = await fetch(`${baseURL}/api/macros`);
    if (macrosResponse.ok) {
      const macros = await macrosResponse.json();
      renderMacrosList(macros);
    }
  } catch (error) {
    // Only log in development
    if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
      console.error('Error loading macros data:', error);
    }
  }
}

/**
 * Render groups list
 */
function renderGroups(groups) {
  const groupsList = document.getElementById('groupsList');

  if (!groups || groups.length === 0) {
    groupsList.innerHTML = '<p class="text-muted">No groups found. Create your first group above.</p>';
    return;
  }

  groupsList.innerHTML = groups.map(group => `
    <div class="card mb-md">
      <div class="card-header">
        <h3 class="card-title">${group.name}</h3>
        <p class="card-subtitle">${group.macro_count || 0} macros</p>
      </div>
      <div class="card-body">
        <div class="macros-in-group" data-group-id="${group.id}">
          <!-- Macros in this group will be listed here -->
        </div>
      </div>
    </div>
  `).join('');
}

/**
 * Render pinned macros
 */
function renderPinnedMacros(pinned) {
  const pinnedList = document.getElementById('pinnedMacrosList');

  if (!pinned || pinned.length === 0) {
    pinnedList.innerHTML = '<p class="text-muted">No pinned macros. Pin macros for quick access.</p>';
    return;
  }

  pinnedList.innerHTML = pinned.map(macro => `
    <div class="card mb-md">
      <div class="card-body">
        <h3>${macro.name}</h3>
        <p>${macro.description || ''}</p>
        <button class="btn btn-ghost btn-sm" onclick="unpinMacro('${macro.id}')">Unpin</button>
      </div>
    </div>
  `).join('');
}

/**
 * Render macros list for manage tab
 */
function renderMacrosList(macros) {
  const macroSelect = document.getElementById('macroSelect');
  const targetGroupSelect = document.getElementById('targetGroupSelect');

  // Populate macro select
  macroSelect.innerHTML = '<option value="">Select a macro...</option>' +
    macros.map(macro => `<option value="${macro.id}">${macro.name}</option>`).join('');

  // Populate target group select (will be populated from groups API)
  // This is a placeholder - actual implementation will fetch groups
}

/**
 * Create new group
 */
document.getElementById('createGroupBtn')?.addEventListener('click', async () => {
  const groupName = document.getElementById('newGroupName').value;
  if (!groupName) {
    alert('Please enter a group name');
    return;
  }

  const baseURL = window.location.origin;
  try {
    const response = await fetch(`${baseURL}/api/macros/groups`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: groupName })
    });

    if (response.ok) {
      document.getElementById('newGroupName').value = '';
      loadMacrosData();
    } else {
      alert('Failed to create group');
    }
  } catch (error) {
    // Only log in development
    if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
      console.error('Error creating group:', error);
    }
    alert('Error creating group');
  }
});

/**
 * Move macro to different group
 */
document.getElementById('moveMacroBtn')?.addEventListener('click', async () => {
  const macroId = document.getElementById('macroSelect').value;
  const targetGroupId = document.getElementById('targetGroupSelect').value;

  if (!macroId || !targetGroupId) {
    alert('Please select both macro and target group');
    return;
  }

  const baseURL = window.location.origin;
  try {
    const response = await fetch(`${baseURL}/api/macros/${macroId}/move`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ group_id: targetGroupId })
    });

    if (response.ok) {
      alert('Macro moved successfully');
      loadMacrosData();
    } else {
      alert('Failed to move macro');
    }
  } catch (error) {
    // Only log in development
    if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
      console.error('Error moving macro:', error);
    }
    alert('Error moving macro');
  }
});

/**
 * Copy macro to different group
 */
document.getElementById('copyMacroBtn')?.addEventListener('click', async () => {
  const macroId = document.getElementById('macroSelect').value;
  const targetGroupId = document.getElementById('targetGroupSelect').value;

  if (!macroId || !targetGroupId) {
    alert('Please select both macro and target group');
    return;
  }

  const baseURL = window.location.origin;
  try {
    const response = await fetch(`${baseURL}/api/macros/${macroId}/copy`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ group_id: targetGroupId })
    });

    if (response.ok) {
      alert('Macro copied successfully');
      loadMacrosData();
    } else {
      alert('Failed to copy macro');
    }
  } catch (error) {
    // Only log in development
    if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
      console.error('Error copying macro:', error);
    }
    alert('Error copying macro');
  }
});

/**
 * Unpin a macro
 */
async function unpinMacro(macroId) {
  const baseURL = window.location.origin;
  try {
    const response = await fetch(`${baseURL}/api/macros/${macroId}/unpin`, {
      method: 'POST'
    });

    if (response.ok) {
      loadMacrosData();
    } else {
      alert('Failed to unpin macro');
    }
  } catch (error) {
    // Only log in development
    if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
      console.error('Error unpinning macro:', error);
    }
    alert('Error unpinning macro');
  }
}

