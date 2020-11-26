const fileInput = document.querySelector('#gif-input input[type=file]');
fileInput.onchange = () => {
  if (fileInput.files.length > 0) {
    const fileName = document.querySelector('#gif-input .file-name');
    fileName.textContent = fileInput.files[0].name;
  }
}
