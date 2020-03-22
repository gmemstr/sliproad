// Register our router, and fire it off initially in case user is being linked a dir.
window.addEventListener("hashchange", router, false);
router()

function getFileListing(provider, path = "") {
  fetch(`/api/files/${provider}${path}`)
    .then((response) => {
      return response.json()
    })
    .then((data) => {
      let files = data["Files"]
      html`
      <div class="grid-sm">
        ${files.map(file =>
          `<a href="${!file.IsDirectory ? `/api/files/${provider}${path}/${file.Name}` : `#${provider + "/" + file.Name}`}">
            ${file.Name}${file.IsDirectory ? '/' : ''}
          </a>
          `
        ).join('')}
      </div>
      `
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
  console.log(path, provider)
  getFileListing(provider, "/" + path)
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