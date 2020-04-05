// Register our router, and fire it off initially in case user is being linked a dir.
window.addEventListener("hashchange", router, false);
router()
let input = ""

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
      html`
      <form action="#" method="post">
        <input type="file" id="file" data-dir="${provider}${path}"><label for="file">Upload</label>
        <progress id="progress" value="0" max="100" hidden=""></progress>
      </form>
      <div class="list">
        ${files.map(file =>
          `<a class="${file.IsDirectory ? "directory" : "file"}" href="${!file.IsDirectory ? `/api/files/${provider}${path}/${file.Name}` : `#${provider}${path === "" ? "" : path}/${file.Name}`}">
            <span>${file.IsDirectory ? '<img src="/icons/folder.svg"/>' : '<img src="/icons/file.svg"/>'}${file.Name}</span>
          </a>
          `
        ).join('')}
      </div>
      `
      // Register our new listeners for uploading files.
      input = document.getElementById("file")
      input.addEventListener("change", onSelectFile, false)
    })
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

// File upload functions. Uses XMLHttpRequest so we can display file upload progress.
function onSelectFile() {
  upload(input.getAttribute("data-dir"), input.files[0])
}
function upload(path, file) {
  let xhrObj = new XMLHttpRequest()
  let formData = new FormData()
  formData.append("file", file)

  xhrObj.upload.addEventListener("loadstart", uploadStarted, false)
  xhrObj.upload.addEventListener("progress", uploadProgress, false)
  xhrObj.upload.addEventListener("load", uploadFinish, false)
  xhrObj.open("POST", `/api/files/${path}`, true)

  xhrObj.send(formData);
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