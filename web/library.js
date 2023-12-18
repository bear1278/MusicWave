document.addEventListener('DOMContentLoaded', function() {
    var token = localStorage.getItem('token');
    var accessToken = localStorage.getItem('accessToken');
    var refreshToken = localStorage.getItem('refreshToken');
    var expiry = localStorage.getItem('expiry');

    const favoritesListContainer = document.getElementById('favorites');
    const myPlaylistListContainer = document.getElementById('my-playlists');
    const addedPlaylistListContainer = document.getElementById('added-playlists');
    const albumListContainer = document.getElementById('album-list');
    const artistListContainer = document.getElementById('artist-list');



    document.getElementById('create-playlist-form').addEventListener('submit', function (event) {
        event.preventDefault();

        var formData = new FormData(event.target);

        fetch('/api/playlists/', {
            method: 'POST',
            headers: {
                'Authorization-1': 'Bearer ' + token,
            },
            body: formData
        })
            .then(response => response.json())
            .then(data => {
                console.log('Ответ от сервера:', data);
            })
            .catch(error => {
                console.error('Ошибка при отправке запроса:', error);
            });
        window.location.assign('/library');
    });




    document.getElementById("create-playlist").addEventListener("click", function () {
        var form = document.getElementById("create-playlist-form");
        if (form.style.display === "none" || form.style.display === "") {
            form.style.display = "flex";
        } else {
            form.style.display = "none";
        }
    });

    fetch('/api/playlists', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization-1': 'Bearer ' + token,
            'Authorization-2': 'Bearer ' + accessToken,
            'Authorization-3': 'Bearer ' + refreshToken,
            'Authorization-4': 'Bearer ' + expiry,
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

            const favorites = data.favorites;

            console.log(favorites)
            const myPlaylists = data.myPlaylists;
            const addedPlaylists = data.addedPlaylists;


            function createFavoritesElement(favorite) {
                const favoriteElement = document.createElement('div');
                favoriteElement.classList.add('track');

                favoriteElement.innerHTML = `
                    <div class="title">
                        <img class="track-icon" src="${favorite.cover}" alt="icon">
                        <div class="names">
                            <p>Name: </p>
                        </div>
                        <div class="track-names">
                            <p><a class="link-names" href="http://localhost:8000/api/playlists/page/${favorite.id}">${favorite.name}</a></p>
                        </div>
                    </div>
                    <div class="time-buttons">
                        <div class="time-info">
                            <img src="/static/images/time.svg" class="timer" alt="duration:">
                            <p class="time">${formatDuration(favorite.duration)}</p>
                        </div>
                       
                    </div>
                `;
                return favoriteElement;
            }
            if (addedPlaylists === null) {
                const noAddedPlaylistsMessage = document.createElement('div');
                noAddedPlaylistsMessage.classList.add('alert');

                const messageParagraph = document.createElement('p');
                messageParagraph.textContent = `You haven't got added playlists.`;
                noAddedPlaylistsMessage.classList.add('text-alert');

                noAddedPlaylistsMessage.appendChild(messageParagraph);
                addedPlaylistListContainer.appendChild(noAddedPlaylistsMessage);
            }else {
                addedPlaylists.forEach(addedPlaylist => {
                    const addedPlaylistElement = createAddedPlaylistElement(addedPlaylist);
                    addedPlaylistListContainer.appendChild(addedPlaylistElement);

                    const button = addedPlaylistElement.querySelector('.button-delete');

                    var dataValue = button.getAttribute('data-value');
                    var jsonData
                    // Проверяем, есть ли значение
                    if (dataValue) {
                        // Парсим JSON из значения data-value
                        button.addEventListener('click', () => {

                            jsonData = JSON.parse(dataValue);
                            handleDeleteButtonLogic(button, jsonData, `api/playlists/exclude/${addedPlaylist.id}`);
                            // Выполняем общую логику, отправляем POST-запрос
                        });
                    } else {
                        console.error('Data-value is missing for the button.');
                    }

                });
            }

                favorites.forEach(favorite => {
                    const favoriteElement = createFavoritesElement(favorite);
                    favoritesListContainer.appendChild(favoriteElement);

                    // const button = myPlaylistElement.querySelector('.button-options');
                    //
                    // var dataValue = button.getAttribute('data-value');
                    // var jsonData
                    // // Проверяем, есть ли значение
                    // if (dataValue) {
                    //     // Парсим JSON из значения data-value
                    //     jsonData = JSON.parse(dataValue);
                    //     handlePostButtonLogic(button, jsonData,`api/playlists/add/${dataValue}`);
                    //     // Выполняем общую логику, отправляем POST-запрос
                    // } else {
                    //     console.error('Data-value is missing for the button.');
                    // }

                });


                function createAddedPlaylistElement(addedPlaylist) {
                    const addedPlaylistElement = document.createElement('div');
                    addedPlaylistElement.classList.add('track');
                    console.log(addedPlaylist)
                    addedPlaylistElement.innerHTML = `
                    <div class="title">
                        <img class="track-icon" src="${addedPlaylist.cover}" alt="icon">
                        <div class="names">
                            <p>Name: </p>
                            <p>Author: </p>
                            <p>Type</p>
                        </div>
                        <div class="track-names">
                            <p><a class="link-names" href="http://localhost:8000/api/playlists/page/${addedPlaylist.id}">${addedPlaylist.name}</a></p>
                            <p>${addedPlaylist.Author.username}</p>
                            <p>${addedPlaylist.type}</p>
                        </div>
                    </div>
                    <div class="time-buttons">
                        <div class="time-info">
                            <img src="/static/images/time.svg" class="timer" alt="duration:">
                            <p class="time">${formatDuration(addedPlaylist.duration)}</p>
                        </div>
                        <button class="button-delete playlist" title="Удалить плейлист из библиотеки" data-value='${addedPlaylist.id}'></button>
                    </div>
                `;
                    return addedPlaylistElement;
                }

                function createMyPlaylistElement(myPlaylist) {
                    const myPlaylistElement = document.createElement('div');
                    myPlaylistElement.classList.add('track');
                    myPlaylistElement.innerHTML = `
                    <div class="title">
                        <img class="track-icon" src="${myPlaylist.cover}" alt="icon">
                        <div class="names">
                            <p>Name: </p>
                            <p>Author: </p>
                            <p>Type: </p>
                        </div>
                        <div class="track-names">
                            <p><a class="link-names" href="http://localhost:8000/api/playlists/page/${myPlaylist.id}">${myPlaylist.name}</a></p>
                            <p>${myPlaylist.Author.username}</p>
                            <p>${myPlaylist.type}</p>
                        </div>
                    </div>
                    <div class="time-buttons">
                        <div class="time-info">
                            <img src="/static/images/time.svg" class="timer" alt="duration:">
                            <p class="time">${formatDuration(myPlaylist.duration)}</p>
                        </div>
                        <div class="buttons">
                            <button class="button-spotify" title="Добавить плейлист в Spotify"></button>
                            <button class="button-delete playlist" title="Удалить плейлист" data-value='${myPlaylist.id}'></button>
                        </div>
                    </div>
                `;
                    return myPlaylistElement;
                }


                if (myPlaylists === null) {
                    const noMyPlaylistsMessage = document.createElement('div');
                    noMyPlaylistsMessage.classList.add('alert');

                    const messageParagraph = document.createElement('p');
                    messageParagraph.textContent = `You haven't got playlists.`;
                    noMyPlaylistsMessage.classList.add('text-alert');

                    noMyPlaylistsMessage.appendChild(messageParagraph);
                    myPlaylistListContainer.appendChild(noMyPlaylistsMessage);
                } else {
                    // Создаем элементы для каждого плейлиста и добавляем их в контейнер
                    myPlaylists.forEach(myPlaylist => {
                        const myPlaylistElement = createMyPlaylistElement(myPlaylist);
                        myPlaylistListContainer.appendChild(myPlaylistElement);

                        const button = myPlaylistElement.querySelector('.button-delete');
                        const buttonSpotify = myPlaylistElement.querySelector('.button-spotify');


                        buttonSpotify.addEventListener('click', () => {
                            // Показываем окно подтверждения
                            const userConfirmed = confirm('Вы хотите добавить плейлист в свой аккаунт спотифай?');

                            // Проверяем ответ пользователя
                            if (userConfirmed) {
                                // Парсим JSON из значения data-value
                                handleAddSpotifyButtonLogic(buttonSpotify, `api/playlists/spotify/${myPlaylist.id}`);

                                // Выполняем общую логику, отправляем POST-запрос
                            } else {
                                console.log('Пользователь отказался от удаления.');
                            }
                        });

                        var dataValue = button.getAttribute('data-value');
                        var jsonData;

                        // Проверяем, есть ли значение
                        if (dataValue) {
                            // Подписываемся на событие click для кнопки
                            button.addEventListener('click', () => {
                                // Показываем окно подтверждения
                                const userConfirmed = confirm('Вы уверены, что хотите удалить этот плейлист?');

                                // Проверяем ответ пользователя
                                if (userConfirmed) {
                                    // Парсим JSON из значения data-value
                                    jsonData = JSON.parse(dataValue);
                                    handleDeleteButtonLogic(button, jsonData, `api/playlists/${myPlaylist.id}`);

                                    // Выполняем общую логику, отправляем POST-запрос
                                } else {
                                    console.log('Пользователь отказался от удаления.');
                                }
                            });
                        } else {
                            console.error('Data-value is missing for the button.');
                        }
                    });
                }



        })

    fetch('/api/albums', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization-1': 'Bearer ' + token,
            'Authorization-2': 'Bearer ' + accessToken,
            'Authorization-3': 'Bearer ' + refreshToken,
            'Authorization-4': 'Bearer ' + expiry,
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

            const albums = data.albums;


            function createAlbumElement(album) {
                const albumElement = document.createElement('div');
                albumElement.classList.add('track');

                const artistNames = album.artists.map(artist => `<a class="link-names" href="http://localhost:8000/api/artists/page/${artist.id}">${artist.name}</a>`).join(', ');

                albumElement.innerHTML = `
                    <div class="title">
                        <img class="track-icon" src="${album.cover}" alt="icon">
                        <div class="names">
                            <p>Name: </p>
                            <p>Artists: </p>
                            <p>Type: </p>
                        </div>
                        <div class="track-names">
                            <p><a class="link-names" href="http://localhost:8000/api/albums/page/${album.id}">${album.name}</a></p>
                            <p>${artistNames}</p>
                            <p>${album.albumType}</p>
                        </div>
                    </div>
                    <div class="time-buttons">
                        <div class="time-info">
                            <img src="/static/images/time.svg" class="timer" alt="duration:">
                            <p class="time">${formatDuration(album.duration)}</p>
                        </div>
                        <div class="buttons">
                            <a class="button-spotify" href="${album.spotifyURL}" target="_blank"></a>
                            <button class="button-delete" title="Удалить альбом из библиотеки" data-value='${JSON.stringify({id: album.id,
                    name: album.name, cover: album.cover, releaseDate: album.releaseDate, albumType: album.albumType,
                    totalTracks: album.totalTracks, spotifyURL: album.spotifyURL, popularity: album.popularity,
                    artists: album.artists})}'></button>
                        </div>
                    </div>
                `;
                return albumElement;
            }

            if (albums === null) {
                const noAlbumsMessage = document.createElement('div');
                noAlbumsMessage.classList.add('alert');

                const messageParagraph = document.createElement('p');
                messageParagraph.textContent = `You haven't got albums.`;
                noAlbumsMessage.classList.add('text-alert');

                noAlbumsMessage.appendChild(messageParagraph);
                albumListContainer.appendChild(noAlbumsMessage);
            }else{
                albums.forEach(album => {
                    console.log(album)
                    const albumElement = createAlbumElement(album);
                    albumListContainer.appendChild(albumElement);

                    const button = albumElement.querySelector('.button-delete');

                    var dataValue = button.getAttribute('data-value');
                    var jsonData
                    // Проверяем, есть ли значение
                    if (dataValue) {
                        // Парсим JSON из значения data-value
                        button.addEventListener('click', () => {

                            jsonData = JSON.parse(dataValue);
                            handleDeleteButtonLogic(button, jsonData, `api/albums/${album.id}`);
                            // Выполняем общую логику, отправляем POST-запрос
                        });
                    } else {
                        console.error('Data-value is missing for the button.');
                    }


                });
            }


        })



    fetch('/api/artists', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization-1': 'Bearer ' + token,
            'Authorization-2': 'Bearer ' + accessToken,
            'Authorization-3': 'Bearer ' + refreshToken,
            'Authorization-4': 'Bearer ' + expiry,
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

            const artists = data.artists;


            function createArtistElement(artist) {
                const artistElement = document.createElement('div');
                artistElement.classList.add('artist');

                artistElement.innerHTML = `
                    <div class="title">
                        <img class="track-icon" src="${artist.picture}" alt="icon">
                        <div class="artist-title">
                            <div class="names">
                                <p>Artist: </p>
                            </div>
                            <div class="artist-names">
                                <p><a class="link-names" href="http://localhost:8000/api/artists/page/${artist.id}">${artist.name}</a></p>
                            </div>
                        </div>
                    </div>
                    <div class="buttons">
                        <a class="button-spotify" href="${artist.spotifyURL}" target="_blank"></a>
                        <button class="button-delete" title="Удалить исполнителя из библиотеки" data-value='${JSON.stringify({
                    id: artist.id, name: artist.name,
                    picture: artist.picture, spotifyURL: artist.spotifyURL, popularity: artist.popularity
                })}'></button>
                    </div>
                `;
                return artistElement;
            }


            if (artists === null) {
                const noArtistsMessage = document.createElement('div');
                noArtistsMessage.classList.add('alert');

                const messageParagraph = document.createElement('p');
                messageParagraph.textContent = `You haven't got albums.`;
                noArtistsMessage.classList.add('text-alert');

                noArtistsMessage.appendChild(messageParagraph);
                artistListContainer.appendChild(noArtistsMessage);
            }else{
                artists.forEach(artist => {
                    console.log(artist)
                    const artistElement = createArtistElement(artist);
                    artistListContainer.appendChild(artistElement);

                    const button = artistElement.querySelector('.button-delete');

                    var dataValue = button.getAttribute('data-value');
                    var jsonData
                    // Проверяем, есть ли значение
                    if (dataValue) {
                        button.addEventListener('click', () => {

                            // Парсим JSON из значения data-value
                            jsonData = JSON.parse(dataValue);
                            handleDeleteButtonLogic(button, jsonData, `api/artists/${artist.id}`);
                            // Выполняем общую логику, отправляем POST-запрос

                        });
                        // Выполняем общую логику, отправляем POST-запрос
                    } else {
                        console.error('Data-value is missing for the button.');
                    }


                });
            }


        })

    function handleDeleteButtonLogic(button, jsonData,url) {
        // Добавляем обработчик события к кнопке для выполнения POST-запроса


        // Проверяем наличие data-value

            // Выполняем POST-запрос
            fetch(url, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization-1': 'Bearer ' + token,
                },
            })
                .then(response => {
                    if (!response.ok) {
                        // Если код ответа не 200, отображаем сообщение в alert
                        return response.json().then(data => {
                            const errorMessage = data.message || 'Unknown error occurred';
                            throw new Error(errorMessage);
                        });
                    }
                    console.log('Track added to playlist successfully');
                })
                .catch(error => {
                    alert(`Error: ${error.message || 'Unknown error occurred'}`);
                    console.error(error);
                });


    }


    function handleAddSpotifyButtonLogic(button, url) {

        fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization-1': 'Bearer ' + token,
                'Authorization-2': 'Bearer ' + accessToken,
                'Authorization-3': 'Bearer ' + refreshToken,
                'Authorization-4': 'Bearer ' + expiry,
            },
        })
            .then(response => {
                if (!response.ok) {
                    // Если код ответа не 200, отображаем сообщение в alert
                    return response.json().then(data => {
                        const errorMessage = data.message || 'Unknown error occurred';
                        throw new Error(errorMessage);
                    });
                }
                console.log('Track added to playlist successfully');
            })
            .catch(error => {
                alert(`Error: ${error.message || 'Unknown error occurred'}`);
                console.error(error);
            });


    }



            function formatDuration(duration) {
        if (duration === 0 || duration === undefined) {
            return `0:00`;
        }
        const minutes = Math.floor(duration / 60000);
        const seconds = ((duration % 60000) / 1000).toFixed(0);
        return `${minutes}:${(seconds < 10 ? '0' : '')}${seconds}`;
    }
});