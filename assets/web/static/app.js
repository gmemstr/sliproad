const $ = selector => document.querySelector(selector);

const init = event => {
    var path = window.location.pathname;
    var pathSplit = path.split("/");
    var path = pathSplit.splice(1).join("/");
    document.title = path + " | PiNAS";

    getFiles(pathSplit[0], path, json => {
        console.log(json);
        document.getElementById("previous-dir").href = "/" + json.Prefix + "/" + json.Previous;
        buildList(json);
    });

    $("#path-header").innerText = path;
    $("#path").value = path;
    document.querySelector("body").addEventListener('contextmenu', function (ev) {
        if (ev.target.localName == "a") {
            ev.preventDefault();

            var d = document.getElementById('context');
            d.classList.remove("hidden");
            d.style.position = "absolute";
            d.style.left = ev.clientX + 'px';
            d.style.top = ev.clientY + 'px';

        }
        return false;
    }, false);

    $("body").addEventListener('click', function (ev) {
        let shouldDismiss = ev.target.dataset.dismissContext == undefined && ev.target.parentElement.classList.contains("context-actions") == false && ev.target.localName != 'a';

        if (ev.which == 1 && shouldDismiss) {
            ev.preventDefault();

            var d = $('#context');
            d.classList.add("hidden");
            return false;
        }
    }, false);
}

const getFiles = (prefix, path, callback) => {
    fetch('/api/' + prefix + '/' + path)
        .then(function (response) {
            return response.json();
        })
        .then(function (jsonResult) {
            callback(jsonResult);
        });
}

const buildList = data => {
    for (var i = 0; i < data.Files.length; i++) {
        let fileItem = document.createElement('p');
        let fileLink = document.createElement('a');
        if (data.Files[i].IsDirectory == true) {
            fileItem.classList.add("directory");
            fileLink.href = "/" + data.Prefix + "/" + data.Path + "/" + data.Files[i].Name;
        } else {
            fileItem.classList.add("file");
            fileLink.href = "/" + data.SinglePrefix + "/" + data.Path + "/" + data.Files[i].Name;
        }

        fileLink.innerText = data.Files[i].Name;
        fileItem.appendChild(fileLink);
        $("#filelist").appendChild(fileItem);
    }
}

const upload = (file, path) => {
    var formData = new FormData();
    formData.append("path", path);
    formData.append("file", file);
    fetch('/upload', { // Your POST endpoint
        method: 'POST',
        body: formData
    }).then(
        success => console.log(success) // Handle the success response object
    ).catch(
        error => console.log(error) // Handle the error response object
    );
};

window.addEventListener('load', init);

