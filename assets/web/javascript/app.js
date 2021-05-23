// Register our router, and fire it off initially in case user is being linked a dir.
window.addEventListener("hashchange", router, false);
router()
let fileinput = ""

// Fetch file listing for a provider and optional path.
function getFileListing(provider, path = "") {
  // There is some funky behaviour happening here between localhost and a deployed instance.
  // This *fixes* is, but it's not ideal.
  if (!path.startsWith("/") && path !== "") {
    path = "/" + path
  }
  fetch(`/api/files/${provider}${path}`)
    .then((response) => {
      return response.json()
    })
    .then((data) => {
      let files = data["Files"]
      if (!files) {
        files = []
      }
      onlyfiles = files.filter(file => !file.IsDirectory)
      onlydirs = files.filter(file => file.IsDirectory)

      html`
      <div class="forms">
      <form id="uploadfile" action="#" method="post">
        <input type="file" id="file" data-dir="${provider}${path}"><label for="file">Upload</label>
        <progress id="progress" value="0" max="100" hidden=""></progress>
      </form>
      <form id="createdir" action="#" method="post">
        <input type="text" id="newdir" data-dir="${provider}${path}">
        <input type="submit" value="Create Directory" id="newdir_submit">
      </form>
      </div>
      <div class="directories list">
        ${onlydirs.map(directory =>
          `<div class="item"><a class="directory" href="#${provider}${path === "" ? "" : path}/${directory.Name}">
            <span>${directory.Name}/</span>
          </a><button onclick="deleteFile('${provider}', '${path === "" ? '' : path}', '${directory.Name}')"><img src="/icons/trash.svg"/></button></div>
          `
        ).join('')}
      </div>

      <div class="files list">
        ${onlyfiles.map(file =>
          `<div class="item"><a class="file" href="/api/files/${provider}${path}/${file.Name}">
            <span><img src="/icons/file.svg"/>${file.Name}</span>
          </a><button onclick="deleteFile('${provider}', '${path === "" ? '' : path}', '${file.Name}')"><img src="/icons/trash.svg"/></button></div>
          `
        ).join('')}
      </div>
      `
      // Register our new listeners for uploading files.
      fileinput = document.getElementById("file")
      fileinput.addEventListener("change", onSelectFile, false)
      createdir = document.getElementById("createdir")
      createdir.addEventListener("submit", mkdir)
    })
}

function deleteFile(provider, path, filename) {
  let xhrObj = new XMLHttpRequest()
  let rp = `${provider}${path === "" ? "" : path}/${filename}`

  xhrObj.addEventListener("loadend", uploadFinish, false)
  xhrObj.open("DELETE", `/api/files/${rp}`, true)

  xhrObj.send()
}

// Fetch list of providers and render.
function getProviders() {
  fetch(`/api/providers`)
    .then((response) => {
      return response.json()
    })
    .then((data) => {
      let providers = data
      html`
        <div class="grid-lg">
        ${providers.map(provider =>
        `<a href=#${provider}>
              ${provider}
        </a>
          `
      ).join('')}
      </div>
      `
    })
}

// Dumb router function for passing around values from the hash.
function router(event = null) {
  let hash = location.hash.replace("#", "")
  // If hash is empty, "redirect" to index.
  if (hash === "") {
    getProviders()
    return
  }

  let path = hash.split("/")
  let provider = path.shift()
  path = path.join("/")
  getFileListing(provider, path)
}

function mkdir(event) {
  event.preventDefault()
  let xhrObj = new XMLHttpRequest()
  mkdir = document.getElementById("newdir")
  let path = mkdir.getAttribute("data-dir")
  let mkdirvalue = mkdir.value

  xhrObj.addEventListener("loadend", uploadFinish, false)
  xhrObj.open("POST", `/api/files/${path}/${mkdirvalue}`, true)
  xhrObj.setRequestHeader("X-NAS-Type", "directory")

  xhrObj.send()
}

// File upload functions. Uses XMLHttpRequest so we can display file upload progress.
function onSelectFile() {
  upload(fileinput.getAttribute("data-dir"), fileinput.files[0])
}
function upload(path, file) {
  let xhrObj = new XMLHttpRequest()
  let formData = new FormData()
  formData.append("file", file)

  xhrObj.upload.addEventListener("loadstart", uploadStarted, false)
  xhrObj.upload.addEventListener("progress", uploadProgress, false)
  xhrObj.upload.addEventListener("loadend", uploadFinish, false)
  xhrObj.open("POST", `/api/files/${path}`, true)

  xhrObj.send(formData)
}

function uploadStarted(e) {
  document.getElementById("progress").hidden = false
}
function uploadProgress(e) {
  let progressBar = document.getElementById("progress")
  progressBar.max = e.total
  progressBar.value = e.loaded
}

function uploadFinish(e) { router() }

// Tagged template function for parsing a string of text as HTML objects
// <3 @innovati for this brilliance.
function html(strings, ...things) {
  // Our "body", where we'll render stuff.
  const body = document.getElementById("main")
  let x = document.createRange().createContextualFragment(
    strings.reduce(
      (markup, string, index) => {
        markup += string

        if (things[index]) {
          markup += things[index]
        }

        return markup
      },
      ''
    )
  )
  body.innerHTML = ""
  body.append(x)
}