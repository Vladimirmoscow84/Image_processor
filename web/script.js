function thumbURL(p) {
    if (!p) return null;
    const normalized = p.startsWith("/") ? p : "/" + p;
    return encodeURI(normalized);
}

async function uploadImage() {
    const file = document.getElementById("imageInput").files[0];
    if (!file) {
        alert("Выберите файл!");
        return;
    }

    const formData = new FormData();
    formData.append("image", file);

    const resp = await fetch("/upload", { method: "POST", body: formData });
    const data = await resp.json();

    alert("Файл отправлен на обработку!");
    loadImages(); 
}

async function loadImages() {
    const listBox = document.getElementById("imageList");
    listBox.innerHTML = "";

    const response = await fetch("/images");
    if (!response.ok) return;

    const images = await response.json();
    images.forEach(img => renderImageCard(img));
}

function renderImageCard(img) {
    const list = document.getElementById("imageList");

    const card = document.createElement("div");
    card.className = "image-card";

    card.innerHTML = `
        <p><b>ID:</b> ${img.id}</p>
        <p><b>Status:</b> ${img.status}</p>
        
        ${img.status === "processed" && img.thumbnailPath
            ? `<img src="${thumbURL(img.thumbnailPath)}" alt="thumb">`
            : `<p>В обработке...</p>`
        }

        <button onclick="deleteImage(${img.id})">Удалить</button>
    `;

    list.appendChild(card);
}

async function deleteImage(id) {
    const resp = await fetch(`/image/${id}`, { method: "DELETE" });
    if (resp.ok) {
        alert("Удалено");
        loadImages();
    }
}


window.addEventListener("load", loadImages);
