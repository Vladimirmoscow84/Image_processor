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

   

    for (let id = 1; id < 50; id++) {
        const response = await fetch(`/image/${id}`);
        if (!response.ok) continue;

        const img = await response.json();
        renderImageCard(img);
    }
}

function renderImageCard(img) {
    const list = document.getElementById("imageList");

    const card = document.createElement("div");
    card.className = "image-card";

    card.innerHTML = `
        <p><b>ID:</b> ${img.id}</p>
        <p><b>Status:</b> ${img.status}</p>
        
        ${img.thumbnailPath 
            ? `<img src="/${img.thumbnailPath}" alt="thumb">`
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
