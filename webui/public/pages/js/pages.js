document.addEventListener("DOMContentLoaded",()=>{loadPages(),initCreatePageButton()});async function loadPages(){const t=window.location.origin,e=document.getElementById("pagesList");try{const n=await fetch(`${t}/api/pages`);if(n.ok){const e=await n.json();renderPages(e)}else e.innerHTML='<p class="text-muted">Unable to load pages</p>'}catch(t){(window.location.hostname==="localhost"||window.location.hostname==="127.0.0.1")&&console.error("Error loading pages:",t),e.innerHTML='<p class="text-muted">Error loading pages</p>'}}function renderPages(e){const t=document.getElementById("pagesList");if(!e||e.length===0){t.innerHTML='<p class="text-muted">No pages found. Create your first page above.</p>';return}t.innerHTML=e.map(e=>`
    <div class="card mb-md">
      <div class="card-body">
        <h3>${e.name||e.title||"Untitled Page"}</h3>
        <p class="text-muted">${e.description||""}</p>
        <div class="form-actions">
          <button class="btn btn-secondary btn-sm" onclick="editPage('${e.id}')">Edit</button>
          <button class="btn btn-outline btn-sm" onclick="viewPage('${e.id}')">View</button>
        </div>
      </div>
    </div>
  `).join("")}function initCreatePageButton(){const e=document.getElementById("createPageBtn");e&&e.addEventListener("click",()=>{showCreatePageModal()})}function showCreatePageModal(){const e=prompt("Enter page name:");e&&createPage(e)}async function createPage(e){const t=window.location.origin;try{const n=await fetch(`${t}/api/pages`,{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify({name:e})});n.ok?loadPages():alert("Failed to create page")}catch(e){(window.location.hostname==="localhost"||window.location.hostname==="127.0.0.1")&&console.error("Error creating page:",e),alert("Error creating page")}}function editPage(e){window.location.href=`/pages/${e}/edit`}function viewPage(e){window.location.href=`/pages/${e}`}