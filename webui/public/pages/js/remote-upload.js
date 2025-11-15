document.addEventListener("DOMContentLoaded",()=>{initUploadForm(),loadUploadHistory()});function initUploadForm(){const e=document.getElementById("uploadForm"),t=document.getElementById("fileInput"),n=document.getElementById("clearBtn");e.addEventListener("submit",handleUpload),n.addEventListener("click",()=>{t.value="",document.getElementById("uploadPath").value="",document.getElementById("uploadResults").innerHTML="",document.getElementById("uploadProgress").style.display="none"})}async function handleUpload(e){e.preventDefault();const i=document.getElementById("fileInput"),o=document.getElementById("uploadPath").value,n=i.files;if(n.length===0){alert("Please select at least one file to upload");return}const a=window.location.origin,s=new FormData;for(let e=0;e<n.length;e++)s.append("files",n[e]);o&&s.append("path",o);const t=document.getElementById("uploadProgress"),r=document.getElementById("progressFill"),c=document.getElementById("progressText");t.style.display="block";try{const e=new XMLHttpRequest;e.upload.addEventListener("progress",e=>{if(e.lengthComputable){const t=e.loaded/e.total*100;r.style.width=`${t}%`,c.textContent=`${Math.round(t)}%`}}),e.addEventListener("load",()=>{if(e.status===200){const t=JSON.parse(e.responseText);showUploadResults(t),loadUploadHistory()}else showUploadError("Upload failed");t.style.display="none"}),e.addEventListener("error",()=>{showUploadError("Network error during upload"),t.style.display="none"}),e.open("POST",`${a}/api/remote-upload`),e.send(s)}catch(e){(window.location.hostname==="localhost"||window.location.hostname==="127.0.0.1")&&console.error("Error uploading files:",e),showUploadError("Error uploading files"),t.style.display="none"}}function showUploadResults(e){const t=document.getElementById("uploadResults");e.success&&e.files?t.innerHTML=`
      <div class="card" style="background-color: rgba(16, 185, 129, 0.1); border-color: #10b981;">
        <div class="card-body">
          <h3>Upload Successful</h3>
          <ul>
            ${e.files.map(e=>`<li>${e.name} - ${e.size||"Uploaded"}</li>`).join("")}
          </ul>
        </div>
      </div>
    `:showUploadError(e.error||"Upload failed")}function showUploadError(e){const t=document.getElementById("uploadResults");t.innerHTML=`
    <div class="card" style="background-color: rgba(239, 68, 68, 0.1); border-color: #ef4444;">
      <div class="card-body">
        <h3>Upload Failed</h3>
        <p>${e}</p>
      </div>
    </div>
  `}async function loadUploadHistory(){const t=window.location.origin,e=document.getElementById("uploadHistory");try{const n=await fetch(`${t}/api/remote-upload/history`);if(n.ok){const e=await n.json();renderUploadHistory(e)}else e.innerHTML='<p class="text-muted">Unable to load upload history</p>'}catch(t){(window.location.hostname==="localhost"||window.location.hostname==="127.0.0.1")&&console.error("Error loading upload history:",t),e.innerHTML='<p class="text-muted">Error loading upload history</p>'}}function renderUploadHistory(e){const t=document.getElementById("uploadHistory");if(!e||e.length===0){t.innerHTML='<p class="text-muted">No upload history available</p>';return}t.innerHTML=e.map(e=>`
    <div class="card mb-md">
      <div class="card-body">
        <h3>${e.filename}</h3>
        <p class="text-muted">Uploaded: ${new Date(e.uploaded_at).toLocaleString()}</p>
        <p class="text-muted">Size: ${formatFileSize(e.size)}</p>
        <p class="text-muted">Path: ${e.path||"/"}</p>
      </div>
    </div>
  `).join("")}function formatFileSize(e){if(e===0)return"0 Bytes";const t=1024,s=["Bytes","KB","MB","GB"],n=Math.floor(Math.log(e)/Math.log(t));return Math.round(e/t**n*100)/100+" "+s[n]}