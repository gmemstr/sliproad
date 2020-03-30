// Register our router, and fire it off initially in case user is being linked a dir.
window.addEventListener("hashchange", router, false);
router()
let input = ""

function getFileListing(provider, path = "") {
  fetch(`/api/files/${provider}${path}`)
    .then((response) => {
      return response.json()
    })
    .then((data) => {
      let files = data["Files"]
      html`
      <form action="#" method="post">
        <input type="file" id="file" data-dir="${provider}${path}"><label for="file">Upload</label>
      </form>
      <div class="grid-sm">
        ${files.map(file =>
          `<a href="${!file.IsDirectory ? `/api/files/${provider}${path}/${file.Name}` : `#${provider}/${path !== "" ? path.replace("/","") + "/" : ""}${file.Name}`}">
            ${file.Name}${file.IsDirectory ? '/' : ''}
          </a>
          `
        ).join('')}
      </div>
      `

      input = document.getElementById("file")
      input.addEventListener("change", onSelectFile, false)
    })
}

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
  getFileListing(provider, "/" + path)
}

function onSelectFile() {
  upload(input.getAttribute("data-dir"), input.files[0])
}
function upload(path, file) {
  let formData = new FormData()
  formData.append("file", file)
  fetch(`/api/files/${path}`, {
    method: "POST",
    body: formData
  }).then(response => response.text())
    .then(text => console.log(text))
    .then(router())
}

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