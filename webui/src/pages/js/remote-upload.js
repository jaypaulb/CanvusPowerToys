/**
 * Remote Upload Page JavaScript
 * Handles file uploads with progress tracking
 */

document.addEventListener('DOMContentLoaded', () => {
  initUploadForm();
  loadUploadHistory();
});

/**
 * Initialize upload form
 */
function initUploadForm() {
  const form = document.getElementById('uploadForm');
  const fileInput = document.getElementById('fileInput');
  const clearBtn = document.getElementById('clearBtn');
  
  form.addEventListener('submit', handleUpload);
  
  clearBtn.addEventListener('click', () => {
    fileInput.value = '';
    document.getElementById('uploadPath').value = '';
    document.getElementById('uploadResults').innerHTML = '';
    document.getElementById('uploadProgress').style.display = 'none';
  });
}

/**
 * Handle file upload
 */
async function handleUpload(e) {
  e.preventDefault();
  
  const fileInput = document.getElementById('fileInput');
  const uploadPath = document.getElementById('uploadPath').value;
  const files = fileInput.files;
  
  if (files.length === 0) {
    alert('Please select at least one file to upload');
    return;
  }
  
  const baseURL = window.location.origin;
  const formData = new FormData();
  
  for (let i = 0; i < files.length; i++) {
    formData.append('files', files[i]);
  }
  
  if (uploadPath) {
    formData.append('path', uploadPath);
  }
  
  // Show progress
  const progressDiv = document.getElementById('uploadProgress');
  const progressFill = document.getElementById('progressFill');
  const progressText = document.getElementById('progressText');
  progressDiv.style.display = 'block';
  
  try {
    const xhr = new XMLHttpRequest();
    
    xhr.upload.addEventListener('progress', (e) => {
      if (e.lengthComputable) {
        const percentComplete = (e.loaded / e.total) * 100;
        progressFill.style.width = `${percentComplete}%`;
        progressText.textContent = `${Math.round(percentComplete)}%`;
      }
    });
    
    xhr.addEventListener('load', () => {
      if (xhr.status === 200) {
        const response = JSON.parse(xhr.responseText);
        showUploadResults(response);
        loadUploadHistory();
      } else {
        showUploadError('Upload failed');
      }
      progressDiv.style.display = 'none';
    });
    
    xhr.addEventListener('error', () => {
      showUploadError('Network error during upload');
      progressDiv.style.display = 'none';
    });
    
    xhr.open('POST', `${baseURL}/api/remote-upload`);
    xhr.send(formData);
    
  } catch (error) {
    console.error('Error uploading files:', error);
    showUploadError('Error uploading files');
    progressDiv.style.display = 'none';
  }
}

/**
 * Show upload results
 */
function showUploadResults(results) {
  const resultsDiv = document.getElementById('uploadResults');
  
  if (results.success && results.files) {
    resultsDiv.innerHTML = `
      <div class="card" style="background-color: rgba(16, 185, 129, 0.1); border-color: #10b981;">
        <div class="card-body">
          <h3>Upload Successful</h3>
          <ul>
            ${results.files.map(file => `<li>${file.name} - ${file.size || 'Uploaded'}</li>`).join('')}
          </ul>
        </div>
      </div>
    `;
  } else {
    showUploadError(results.error || 'Upload failed');
  }
}

/**
 * Show upload error
 */
function showUploadError(message) {
  const resultsDiv = document.getElementById('uploadResults');
  resultsDiv.innerHTML = `
    <div class="card" style="background-color: rgba(239, 68, 68, 0.1); border-color: #ef4444;">
      <div class="card-body">
        <h3>Upload Failed</h3>
        <p>${message}</p>
      </div>
    </div>
  `;
}

/**
 * Load upload history
 */
async function loadUploadHistory() {
  const baseURL = window.location.origin;
  const historyDiv = document.getElementById('uploadHistory');
  
  try {
    const response = await fetch(`${baseURL}/api/remote-upload/history`);
    if (response.ok) {
      const history = await response.json();
      renderUploadHistory(history);
    } else {
      historyDiv.innerHTML = '<p class="text-muted">Unable to load upload history</p>';
    }
  } catch (error) {
    console.error('Error loading upload history:', error);
    historyDiv.innerHTML = '<p class="text-muted">Error loading upload history</p>';
  }
}

/**
 * Render upload history
 */
function renderUploadHistory(history) {
  const historyDiv = document.getElementById('uploadHistory');
  
  if (!history || history.length === 0) {
    historyDiv.innerHTML = '<p class="text-muted">No upload history available</p>';
    return;
  }
  
  historyDiv.innerHTML = history.map(item => `
    <div class="card mb-md">
      <div class="card-body">
        <h3>${item.filename}</h3>
        <p class="text-muted">Uploaded: ${new Date(item.uploaded_at).toLocaleString()}</p>
        <p class="text-muted">Size: ${formatFileSize(item.size)}</p>
        <p class="text-muted">Path: ${item.path || '/'}</p>
      </div>
    </div>
  `).join('');
}

/**
 * Format file size
 */
function formatFileSize(bytes) {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}

