function dragEnter(e) {
    document.getElementById("uploadoverlay").classList.remove("hidden");
}

function dragExit(e) {
    document.getElementById("uploadoverlay").classList.add("hidden");
}
document.addEventListener("drop", function(e) {
    var files = e.target.files || e.dataTransfer.files
    console.log(files)
    document.getElementById("uploadoverlay").classList.add("hidden");
    for (var i = 0; i <= files.length; i++) {
        upload(files[i]);
    }
})

window.addEventListener("dragover",function(e){
    e = e || event;
    e.preventDefault();
},false);
window.addEventListener("drop",function(e){
    e = e || event;
    e.preventDefault();
},false);


var upload = (file) => {
    var path = document.getElementById("path").value;
    var formData  = new FormData();
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