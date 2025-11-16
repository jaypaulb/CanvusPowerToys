/**
 * Macros Page JavaScript
 * Handles widget management: move, copy, grouping, and pinning
 */

document.addEventListener("DOMContentLoaded", () => {
  console.log("[macros.js] DOMContentLoaded - Initializing macros page");

  // 1) Setup tab navigation
  console.log("[macros.js] Setting up tabs");
  setupTabs();

  // 2) Fetch zones to populate dropdowns
  console.log("[macros.js] Fetching zones");
  fetchZones();

  // 3) Bind button clicks for Manage (Move/Copy)
  const moveButton = document.getElementById("moveButton");
  const copyButton = document.getElementById("copyButton");

  console.log("[macros.js] Binding Manage buttons:", { moveButton: !!moveButton, copyButton: !!copyButton });
  if (moveButton) {
    moveButton.addEventListener("click", () => {
      console.log("[macros.js] Move button clicked");
      manageMove();
    });
  } else {
    console.error("[macros.js] ERROR: moveButton not found!");
  }

  if (copyButton) {
    copyButton.addEventListener("click", () => {
      console.log("[macros.js] Copy button clicked");
      manageCopy();
    });
  } else {
    console.error("[macros.js] ERROR: copyButton not found!");
  }

  // 4) Bind Grouping logic
  const autoGridButton = document.getElementById("autoGridButton");
  const groupColorButton = document.getElementById("groupColorButton");
  const groupTitleButton = document.getElementById("groupTitleButton");

  console.log("[macros.js] Binding Grouping buttons:", {
    autoGridButton: !!autoGridButton,
    groupColorButton: !!groupColorButton,
    groupTitleButton: !!groupTitleButton
  });
  if (autoGridButton) {
    autoGridButton.addEventListener("click", () => {
      console.log("[macros.js] Auto Grid button clicked");
      autoGrid();
    });
  } else {
    console.error("[macros.js] ERROR: autoGridButton not found!");
  }

  if (groupColorButton) {
    groupColorButton.addEventListener("click", () => {
      console.log("[macros.js] Group by Color button clicked");
      groupByColor();
    });
  } else {
    console.error("[macros.js] ERROR: groupColorButton not found!");
  }

  if (groupTitleButton) {
    groupTitleButton.addEventListener("click", () => {
      console.log("[macros.js] Group by Title button clicked");
      groupByTitle();
    });
  } else {
    console.error("[macros.js] ERROR: groupTitleButton not found!");
  }

  // 5) Bind Pinning logic
  const pinAllButton = document.getElementById("pinAllButton");
  const unpinAllButton = document.getElementById("unpinAllButton");

  console.log("[macros.js] Binding Pinning buttons:", { pinAllButton: !!pinAllButton, unpinAllButton: !!unpinAllButton });
  if (pinAllButton) {
    pinAllButton.addEventListener("click", () => {
      console.log("[macros.js] Pin All button clicked");
      pinAll();
    });
  } else {
    console.error("[macros.js] ERROR: pinAllButton not found!");
  }

  if (unpinAllButton) {
    unpinAllButton.addEventListener("click", () => {
      console.log("[macros.js] Unpin All button clicked");
      unpinAll();
    });
  } else {
    console.error("[macros.js] ERROR: unpinAllButton not found!");
  }

  // 6) Color tolerance slider
  console.log("[macros.js] Setting up color tolerance slider");
  setupColorToleranceSlider();

  console.log("[macros.js] Initialization complete");
});

/* ------------------------------ TAB SWITCHING ------------------------------ */
function setupTabs() {
  const tabButtons = document.querySelectorAll('.tab-button');
  const tabContents = document.querySelectorAll('.tab-content');

  tabButtons.forEach(button => {
    button.addEventListener('click', () => {
      // Remove active class from all buttons and contents
      tabButtons.forEach(btn => btn.classList.remove('active'));
      tabContents.forEach(content => content.classList.remove('active'));

      // Add active class to clicked button and corresponding content
      button.classList.add('active');
      const tabId = button.getAttribute('data-tab');
      const content = document.getElementById(`${tabId}-content`);
      if (content) {
        content.classList.add('active');
      }
    });
  });
}

/* ------------------------------ FETCH ZONES ------------------------------ */
async function fetchZones() {
  try {
    console.log("[fetchZones] Fetching zones and canvas details...");
    const res = await fetch("/get-zones", {
      headers: {
        'Cache-Control': 'no-cache'
      }
    });
    const data = await res.json();

    if (!data.success || !data.zones) {
      throw new Error("Failed to retrieve zones from the server.");
    }

    console.log(`[fetchZones] Retrieved ${data.zones.length} zones.`);

    // Populate the dropdowns with the updated zones list
    populateZoneDropdowns(data.zones);
  } catch (err) {
    console.error("[fetchZones] Error:", err.message);
    displayMessage(err.message, "error");
  }
}

function populateZoneDropdowns(zones) {
  try {
    const dropdowns = {
      "manageSourceZone": document.getElementById("manageSourceZone"),
      "manageTargetZone": document.getElementById("manageTargetZone"),
      "arrangeSourceZone": document.getElementById("arrangeSourceZone"),
      "pinSourceZone": document.getElementById("pinSourceZone")
    };

    // Clear existing options and add default
    Object.values(dropdowns).forEach(dropdown => {
      if (dropdown) {
        dropdown.innerHTML = '<option value="">Select a zone...</option>';
      }
    });

    // Sort zones by anchor_name
    const sortedZones = [...zones].sort((a, b) => {
      const nameA = (a.anchor_name || '').toLowerCase();
      const nameB = (b.anchor_name || '').toLowerCase();
      return nameA.localeCompare(nameB, undefined, { numeric: true });
    });

    // Add zone options
    sortedZones.forEach((zone) => {
      const zoneName = zone.anchor_name || `Zone ${zone.id}`;
      const option = document.createElement("option");
      option.value = zone.id;
      option.textContent = zoneName;

      // Add to each dropdown
      Object.values(dropdowns).forEach(dropdown => {
        if (dropdown) {
          dropdown.appendChild(option.cloneNode(true));
        }
      });
    });
  } catch (err) {
    console.error("[populateZoneDropdowns] Error:", err.message);
    displayMessage("Error populating zone dropdowns: " + err.message, "error");
  }
}

/* ------------------------------ MANAGE ACTIONS ------------------------------ */
async function manageMove() {
  console.log("[macros.js] manageMove() called");
  const sourceZoneId = document.getElementById("manageSourceZone")?.value;
  const targetZoneId = document.getElementById("manageTargetZone")?.value;
  console.log("[macros.js] Zone IDs:", { sourceZoneId, targetZoneId });

  if (!sourceZoneId || !targetZoneId) {
    const msg = "Please select both Source and Target zones.";
    console.warn("[macros.js] Validation failed:", msg);
    displayMessage(msg, "error");
    return;
  }
  try {
    const payload = { sourceZoneId, targetZoneId };
    console.log("[macros.js] Sending POST /api/macros/move with payload:", payload);
    const resp = await postJson("/api/macros/move", payload);
    console.log("[macros.js] Move response:", resp);
    displayMessage(resp.message || "Widgets moved successfully", "success");
  } catch (err) {
    console.error("[macros.js] Move failed:", err);
    displayMessage(err.message || "Failed to move widgets", "error");
  }
}

async function manageCopy() {
  console.log("[macros.js] manageCopy() called");
  const sourceZoneId = document.getElementById("manageSourceZone")?.value;
  const targetZoneId = document.getElementById("manageTargetZone")?.value;
  console.log("[macros.js] Zone IDs:", { sourceZoneId, targetZoneId });

  if (!sourceZoneId || !targetZoneId) {
    const msg = "Please select both Source and Target zones.";
    console.warn("[macros.js] Validation failed:", msg);
    displayMessage(msg, "error");
    return;
  }
  try {
    const payload = { sourceZoneId, targetZoneId };
    console.log("[macros.js] Sending POST /api/macros/copy with payload:", payload);
    const resp = await postJson("/api/macros/copy", payload);
    console.log("[macros.js] Copy response:", resp);
    displayMessage(resp.message || "Widgets copied successfully", "success");
  } catch (err) {
    console.error("[macros.js] Copy failed:", err);
    displayMessage(err.message || "Failed to copy widgets", "error");
  }
}

/* ------------------------------ ARRANGE ACTIONS ------------------------------ */
async function autoGrid() {
  console.log("[macros.js] autoGrid() called");
  const sourceZoneId = document.getElementById("arrangeSourceZone")?.value;
  console.log("[macros.js] Zone ID:", sourceZoneId);

  if (!sourceZoneId) {
    const msg = "Please select a Source zone.";
    console.warn("[macros.js] Validation failed:", msg);
    displayMessage(msg, "error");
    return;
  }
  try {
    const payload = { zoneId: sourceZoneId };
    console.log("[macros.js] Sending POST /api/macros/auto-grid with payload:", payload);
    const resp = await postJson("/api/macros/auto-grid", payload);
    console.log("[macros.js] Auto grid response:", resp);
    displayMessage(resp.message || "Auto grid applied successfully", "success");
  } catch (err) {
    console.error("[macros.js] Auto grid failed:", err);
    displayMessage(err.message || "Failed to apply auto grid", "error");
  }
}

async function groupByColor() {
  console.log("[macros.js] groupByColor() called");
  const sourceZoneId = document.getElementById("arrangeSourceZone")?.value;
  const colorTolerance = document.getElementById("colorToleranceSlider")?.value;
  console.log("[macros.js] Zone ID:", sourceZoneId, "Color tolerance:", colorTolerance);

  if (!sourceZoneId) {
    const msg = "Please select a Source zone.";
    console.warn("[macros.js] Validation failed:", msg);
    displayMessage(msg, "error");
    return;
  }
  try {
    const payload = { zoneId: sourceZoneId, colorTolerance: parseInt(colorTolerance) };
    console.log("[macros.js] Sending POST /api/macros/group-color with payload:", payload);
    const resp = await postJson("/api/macros/group-color", payload);
    console.log("[macros.js] Group by color response:", resp);
    displayMessage(resp.message || "Grouped by color successfully", "success");
  } catch (err) {
    console.error("[macros.js] Group by color failed:", err);
    displayMessage(err.message || "Failed to group by color", "error");
  }
}

async function groupByTitle() {
  console.log("[macros.js] groupByTitle() called");
  const sourceZoneId = document.getElementById("arrangeSourceZone")?.value;
  console.log("[macros.js] Zone ID:", sourceZoneId);

  if (!sourceZoneId) {
    const msg = "Please select a Source zone.";
    console.warn("[macros.js] Validation failed:", msg);
    displayMessage(msg, "error");
    return;
  }
  try {
    const payload = { zoneId: sourceZoneId };
    console.log("[macros.js] Sending POST /api/macros/group-title with payload:", payload);
    const resp = await postJson("/api/macros/group-title", payload);
    console.log("[macros.js] Group by title response:", resp);
    displayMessage(resp.message || "Grouped by title successfully", "success");
  } catch (err) {
    console.error("[macros.js] Group by title failed:", err);
    displayMessage(err.message || "Failed to group by title", "error");
  }
}

/* ------------------------------ PIN ACTIONS ------------------------------ */
async function pinAll() {
  console.log("[macros.js] pinAll() called");
  const sourceZoneId = document.getElementById("pinSourceZone")?.value;
  console.log("[macros.js] Zone ID:", sourceZoneId);

  if (!sourceZoneId) {
    const msg = "Please select a Source zone.";
    console.warn("[macros.js] Validation failed:", msg);
    displayMessage(msg, "error");
    return;
  }
  try {
    const payload = { zoneId: sourceZoneId };
    console.log("[macros.js] Sending POST /api/macros/pin-all with payload:", payload);
    const resp = await postJson("/api/macros/pin-all", payload);
    console.log("[macros.js] Pin all response:", resp);
    displayMessage(resp.message || "All widgets pinned successfully", "success");
  } catch (err) {
    console.error("[macros.js] Pin all failed:", err);
    displayMessage(err.message || "Failed to pin widgets", "error");
  }
}

async function unpinAll() {
  console.log("[macros.js] unpinAll() called");
  const sourceZoneId = document.getElementById("pinSourceZone")?.value;
  console.log("[macros.js] Zone ID:", sourceZoneId);

  if (!sourceZoneId) {
    const msg = "Please select a Source zone.";
    console.warn("[macros.js] Validation failed:", msg);
    displayMessage(msg, "error");
    return;
  }
  try {
    const payload = { zoneId: sourceZoneId };
    console.log("[macros.js] Sending POST /api/macros/unpin-all with payload:", payload);
    const resp = await postJson("/api/macros/unpin-all", payload);
    console.log("[macros.js] Unpin all response:", resp);
    displayMessage(resp.message || "All widgets unpinned successfully", "success");
  } catch (err) {
    console.error("[macros.js] Unpin all failed:", err);
    displayMessage(err.message || "Failed to unpin widgets", "error");
  }
}

/* ------------------------------ COLOR TOLERANCE SLIDER ------------------------------ */
function setupColorToleranceSlider() {
  const slider = document.getElementById("colorToleranceSlider");
  const valueDisplay = document.getElementById("colorToleranceValue");

  if (slider && valueDisplay) {
    slider.addEventListener("input", (e) => {
      valueDisplay.textContent = e.target.value + "%";
    });
  }
}

/* ------------------------------ HELPER FUNCTIONS ------------------------------ */
async function postJson(url, payload) {
  console.log("[macros.js] postJson() - URL:", url, "Payload:", payload);
  const response = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(payload)
  });

  console.log("[macros.js] postJson() - Response status:", response.status, response.statusText);

  if (!response.ok) {
    const errorText = await response.text();
    console.error("[macros.js] postJson() - Error response:", errorText);
    let error;
    try {
      error = JSON.parse(errorText);
    } catch {
      error = { error: errorText || "Request failed" };
    }
    throw new Error(error.error || `HTTP ${response.status}`);
  }

  const result = await response.json();
  console.log("[macros.js] postJson() - Success response:", result);
  return result;
}

function displayMessage(text, type) {
  const messageEl = document.getElementById("manageMessage") ||
                    document.getElementById("arrangeMessage") ||
                    document.getElementById("pinMessage");

  if (messageEl) {
    messageEl.textContent = text;
    messageEl.className = `message ${type} mt-md`;
    messageEl.style.display = "block";

    // Auto-hide after 5 seconds
    setTimeout(() => {
      messageEl.style.display = "none";
    }, 5000);
  } else {
    console.log(`[${type}] ${text}`);
  }
}
