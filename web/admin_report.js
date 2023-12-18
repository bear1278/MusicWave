document.getElementById('report').addEventListener("submit", function (event) {
    var token = localStorage.getItem('token');
    event.preventDefault();

    var selectedRadio = document.querySelector('input[name="format"]:checked');

    if (selectedRadio) {

        var url='/admin/history/'+selectedRadio.value

        // Construct the URL
        SendRequest(url,token,selectedRadio.value)
    } else {
        alert('Please enter search text and select a search type.');
    }


});

function SendRequest(url,token,fileType){
    fetch(url, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization-1': 'Bearer ' + token,

        },
    })
        .then(response =>
        {
            if (!response.ok){
                throw Error(response.json().message)
            }
            return response.blob()
        })
        .then(blob => {
            // Запрашиваем у пользователя выбор директории для сохранения файла
            if (fileType==='pdf'){
                saveAs(blob, 'report.pdf');
            }else{
                saveAs(blob, 'report.xlsx');
            }

        })
        .catch(error => {
            console.error('Ошибка при скачивании файла:', error);
        });
}

document.getElementById('export').addEventListener("click",function (){
    var token = localStorage.getItem('token');
    fetch('/admin/export', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization-1': 'Bearer ' + token,

        },
    })
        .then(response =>
        {
            if (!response.ok){
                throw Error(response.json().message)
            }
            return response.blob()
        })
        .then(blob => {
            // Запрашиваем у пользователя выбор директории для сохранения файла

            saveAs(blob, 'export.json');

        })
        .catch(error => {
            console.error('Ошибка при скачивании файла:', error);
        });
})

document.getElementById('import').addEventListener('submit', function (event) {
    event.preventDefault();

    var formData = new FormData(event.target);

    var token = localStorage.getItem('token');

    fetch('/admin/import', {
        method: 'POST',
        headers: {
            'Authorization-1': 'Bearer ' + token,
        },
        body: formData
    })
        .then(response => {
            if(response.ok){
                alert('Success.');
            }
            else{

                // Если код ответа не 200, отображаем сообщение в alert
                return response.json().then(data => {
                    const errorMessage = data.message || 'Unknown error occurred';
                    throw new Error(errorMessage);
                });

            }
        })
        .catch(error => {
            console.error('Ошибка при отправке запроса:', error);
            alert('Ошибка: ' + error.message);
        });


});


document.addEventListener('DOMContentLoaded', function() {
    var token = localStorage.getItem('token');


    const trackListContainer = document.getElementById('track-list');


    fetch('/admin/history', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization-1': 'Bearer ' + token,
        },
    })
        .then(response => {
            if (!response.ok) {
                if (response.status === 401) {
                    alert('You need to enter to your account.');
                    window.location.assign('/auth/sign-in');
                }
                throw new Error(`HTTP error! Status: ${response.status} - ${response.statusText}`);
            }
            return response.json();
        })
        .then(data => {
            const history = data.history;


            // Используем функцию для создания элемента трека
            function createTrackElement(track) {
                const trackElement = document.createElement('div');
                trackElement.classList.add('track');


                trackElement.innerHTML = `
                    <div class="title">
                        <p class="text">${track.id}</p>
                        <div class="names">
                            <p>Name: </p>
                            <p>Type: </p>
                            <p>Reason: </p>
                        </div>
                        <div class="track-names">
                            <p class="link-names">${track.name}</p>
                            <p class="link-names">${track.type}</p>
                            <p class="link-names" >${track.reason}</p>
                        </div>
                        <div class="names">
                            <p>Deleted Date: </p>
                        </div>
                        <div class="track-names">
                            <p class="link-names">${track.deletedDate}</p>
                        </div>
                    </div>
                    
                `;
                return trackElement;
            }

            if (history === null) {
                const noAlbumsMessage = document.createElement('div');
                noAlbumsMessage.classList.add('alert');

                const messageParagraph = document.createElement('p');
                messageParagraph.textContent = `We haven't albums to suggest you.`;
                noAlbumsMessage.classList.add('text-alert');

                noAlbumsMessage.appendChild(messageParagraph);
                albumListContainer.appendChild(noAlbumsMessage);
            }else{
                // Создаем элементы для каждого альбома и добавляем их в контейнер
                history.forEach(record => {
                    const Element = createTrackElement(record);
                    trackListContainer.appendChild(Element);

                });
            }


        });
})
