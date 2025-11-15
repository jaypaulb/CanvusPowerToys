/**
 * Pages Management JavaScript
 * Handles page listing, creation, and editing
 */

document.addEventListener('DOMContentLoaded', () => {
  loadPages();
  initCreatePageButton();
});

/**
 * Load pages from API
 */
async function loadPages() {
  const baseURL = window.location.origin;
  const pagesList = document.getElementById('pagesList');
  
  try {
    const response = await fetch(`${baseURL}/api/pages`);
    if (response.ok) {
      const pages = await response.json();
      renderPages(pages);
    } else {
      pagesList.innerHTML = '<p class="text-muted">Unable to load pages</p>';
    }
  } catch (error) {
    console.error('Error loading pages:', error);
    pagesList.innerHTML = '<p class="text-muted">Error loading pages</p>';
  }
}

/**
 * Render pages list
 */
function renderPages(pages) {
  const pagesList = document.getElementById('pagesList');
  
  if (!pages || pages.length === 0) {
    pagesList.innerHTML = '<p class="text-muted">No pages found. Create your first page above.</p>';
    return;
  }
  
  pagesList.innerHTML = pages.map(page => `
    <div class="card mb-md">
      <div class="card-body">
        <h3>${page.name || page.title || 'Untitled Page'}</h3>
        <p class="text-muted">${page.description || ''}</p>
        <div class="form-actions">
          <button class="btn btn-secondary btn-sm" onclick="editPage('${page.id}')">Edit</button>
          <button class="btn btn-outline btn-sm" onclick="viewPage('${page.id}')">View</button>
        </div>
      </div>
    </div>
  `).join('');
}

/**
 * Initialize create page button
 */
function initCreatePageButton() {
  const createBtn = document.getElementById('createPageBtn');
  if (createBtn) {
    createBtn.addEventListener('click', () => {
      showCreatePageModal();
    });
  }
}

/**
 * Show create page modal
 */
function showCreatePageModal() {
  // This would use the modal template
  // For now, simple prompt
  const pageName = prompt('Enter page name:');
  if (pageName) {
    createPage(pageName);
  }
}

/**
 * Create new page
 */
async function createPage(name) {
  const baseURL = window.location.origin;
  
  try {
    const response = await fetch(`${baseURL}/api/pages`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name })
    });
    
    if (response.ok) {
      loadPages();
    } else {
      alert('Failed to create page');
    }
  } catch (error) {
    console.error('Error creating page:', error);
    alert('Error creating page');
  }
}

/**
 * Edit page
 */
function editPage(pageId) {
  // Navigate to edit page or show edit modal
  window.location.href = `/pages/${pageId}/edit`;
}

/**
 * View page
 */
function viewPage(pageId) {
  // Navigate to view page
  window.location.href = `/pages/${pageId}`;
}

