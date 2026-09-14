const API_BASE = "http://localhost:8080";

/* 
STATE
 */

let currentGame = null;
let currentGameId = null;
let authMode = "login";

/* 
DOM
 */

const authSection = document.getElementById("authSection");
const menuSection = document.getElementById("menuSection");
const gameSection = document.getElementById("gameSection");

const loginTab = document.getElementById("loginTab");
const signupTab = document.getElementById("signupTab");

const authForm = document.getElementById("authForm");

const loginInput = document.getElementById("loginInput");
const passwordInput = document.getElementById("passwordInput");

const authButton = document.getElementById("authButton");
const authStatus = document.getElementById("authStatus");

const logoutBtn = document.getElementById("logoutBtn");
const userInfo = document.getElementById("userInfo");

const computerGameBtn = document.getElementById("computerGameBtn");
const playerGameBtn = document.getElementById("playerGameBtn");
const refreshGamesBtn = document.getElementById("refreshGamesBtn");
const historyBtn = document.getElementById("historyBtn");
const leaderboardBtn = document.getElementById("leaderboardBtn");

const extraSection = document.getElementById("extraSection");
const extraTitle = document.getElementById("extraTitle");
const extraContent = document.getElementById("extraContent");

const gamesList = document.getElementById("gamesList");
const menuStatus = document.getElementById("menuStatus");

const gameIdElement = document.getElementById("gameId");
const gameStatus = document.getElementById("gameStatus");
const board = document.getElementById("board");
const backBtn = document.getElementById("backBtn");
const refreshGameBtn = document.getElementById("refreshGameBtn");

/* 
AUTH STORAGE
 */

function getTokens() {
    const accessToken = localStorage.getItem("accessToken");
    const refreshToken = localStorage.getItem("refreshToken");

    if (!accessToken || !refreshToken) {
        return null;
    }

    return { accessToken, refreshToken };
}

function saveTokens(accessToken, refreshToken) {
    if (accessToken) localStorage.setItem("accessToken", accessToken);
    if (refreshToken) localStorage.setItem("refreshToken", refreshToken);
}

function clearCredentials() {
    localStorage.removeItem("accessToken");
    localStorage.removeItem("refreshToken");
    localStorage.removeItem("userUUID");
    localStorage.removeItem("userLogin");
}

function getAuthHeader() {
    const tokens = getTokens();
    if (!tokens) return null;
    return `Bearer ${tokens.accessToken}`;
}

function getAuthHeaders() {
    const headers = {
        "Content-Type": "application/json"
    };
    const authorization = getAuthHeader();
    if (authorization) {
        headers.Authorization = authorization;
    }
    return headers;
}

async function refreshAccessToken() {
    const tokens = getTokens();

    if (!tokens || !tokens.refreshToken) {
        throw new Error("Refresh token required");
    }

    const response = await fetch(`${API_BASE}/refreshAccessToken`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            refreshToken: tokens.refreshToken,
            refresh_token: tokens.refreshToken
        })
    });

    const data = await parseResponse(response);

    if (!response.ok) {
        clearCredentials();
        throw new Error(getErrorMessage(data, "Не удалось обновить токен"));
    }

    const newAccess = data.accessToken || data.access_token;
    const newRefresh = data.refreshToken || data.refresh_token || tokens.refreshToken;

    if (newAccess) {
        saveTokens(newAccess, newRefresh);
        return { accessToken: newAccess, refreshToken: newRefresh };
    }

    clearCredentials();
    throw new Error("Неверный формат ответа сервера при обновлении токена");
}

/* 
API HELPERS
 */

async function parseResponse(response) {
    const text = await response.text();

    console.group("HTTP response", response.status, response.statusText, response.url);
    console.log("headers:", [...response.headers.entries()]);
    console.log("body text:", text);
    console.groupEnd();

    if (!text) return null;

    try {
        return JSON.parse(text);
    } catch {
        return text;
    }
}

async function fetchWithAuth(url, options = {}) {
    options.headers = {
        ...getAuthHeaders(),
        ...(options.headers || {})
    };

    console.group("HTTP request", options.method || "GET", url);
    console.log("options:", options);
    console.groupEnd();

    let response = await fetch(url, options);

    if (response.status === 401) {
        console.warn("401 received, refreshing access token");
        try {
            await refreshAccessToken();

            options.headers = {
                ...options.headers,
                ...getAuthHeaders()
            };

            console.log("Retrying request after refresh", options.method || "GET", url);
            response = await fetch(url, options);
        } catch (error) {
            console.error("Token refresh failed", error);
            handleUnauthorized();
            return response;
        }
    }

    return response;
}

function getErrorMessage(data, defaultMessage) {
    if (!data) return defaultMessage;
    if (typeof data === "string") return data;
    return data.message || data.error || defaultMessage;
}

/* 
UI
 */

function showSection(section) {
    authSection.classList.add("hidden");
    menuSection.classList.add("hidden");
    gameSection.classList.add("hidden");

    section.classList.remove("hidden");
}

function showAuth() {
    showSection(authSection);
    loginInput.focus();
}

function showMenu() {
    showSection(menuSection);
    hideExtraSection();

    const userLogin = localStorage.getItem("userLogin");
    if (userLogin) {
        userInfo.textContent = `Пользователь: ${userLogin}`;
    }

    loadAvailableGames();
}

function showExtraSection(title, html) {
    extraTitle.textContent = title;
    extraContent.innerHTML = html;
    extraSection.classList.remove("hidden");
}

function hideExtraSection() {
    if (extraSection) {
        extraSection.classList.add("hidden");
        extraTitle.textContent = "";
        extraContent.innerHTML = "";
    }
}

function showGame() {
    showSection(gameSection);
}

function setAuthStatus(message, type = "info") {
    authStatus.textContent = message;
    authStatus.className = `status ${type}`;
}

function setMenuStatus(message, type = "info") {
    menuStatus.textContent = message;
    menuStatus.className = `status ${type}`;
}

/* 
AUTH TABS
 */

function setAuthMode(mode) {
    authMode = mode;

    loginTab.classList.remove("active");
    signupTab.classList.remove("active");

    if (mode === "login") {
        loginTab.classList.add("active");
        authButton.textContent = "Войти";
        passwordInput.autocomplete = "current-password";
    } else {
        signupTab.classList.add("active");
        authButton.textContent = "Зарегистрироваться";
        passwordInput.autocomplete = "new-password";
    }

    setAuthStatus("");
}

loginTab.addEventListener("click", () => setAuthMode("login"));
signupTab.addEventListener("click", () => setAuthMode("signup"));

/* 
SIGN UP
 */

async function register() {
    const login = loginInput.value.trim();
    const password = passwordInput.value;

    if (!login) {
        setAuthStatus("Введите логин", "error");
        loginInput.focus();
        return;
    }

    if (!password) {
        setAuthStatus("Введите пароль", "error");
        passwordInput.focus();
        return;
    }

    try {
        setAuthStatus("Регистрация...", "info");
        authButton.disabled = true;

        const response = await fetch(`${API_BASE}/auth/register`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ login, password })
        });

        const data = await parseResponse(response);

        if (!response.ok) {
            throw new Error(getErrorMessage(data, "Не удалось зарегистрироваться"));
        }

        setAuthStatus("Регистрация успешна. Выполняем вход...", "success");
        await loginUser();

    } catch (error) {
        console.error("Registration error:", error);
        clearCredentials();
        setAuthStatus(error.message, "error");
    } finally {
        authButton.disabled = false;
    }
}

/* 
LOGIN
 */

async function loginUser() {
    const login = loginInput.value.trim();
    const password = passwordInput.value;

    if (!login || !password) {
        setAuthStatus("Введите логин и пароль", "error");
        return;
    }

    try {
        setAuthStatus("Авторизация...", "info");
        authButton.disabled = true;

        const response = await fetch(`${API_BASE}/auth/login`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ login, password })
        });

        const data = await parseResponse(response);

        if (!response.ok) {
            throw new Error(getErrorMessage(data, "Неверный логин или пароль"));
        }

        if (data && typeof data === "object") {
            const accessToken = data.accessToken || data.access_token;
            const refreshToken = data.refreshToken || data.refresh_token;

            if (accessToken && refreshToken) {
                saveTokens(accessToken, refreshToken);
            }

            if (data.login) {
                localStorage.setItem("userLogin", data.login);
            }

            const uuid = data.uuid || data.UUID || data.id;
            if (uuid) {
                localStorage.setItem("userUUID", uuid);
            }
        }

        setAuthStatus("", "info");
        showMenu();

    } catch (error) {
        console.error("Login error:", error);
        clearCredentials();
        setAuthStatus(error.message, "error");
    } finally {
        authButton.disabled = false;
    }
}

/* 
AUTH FORM
 */

authForm.addEventListener("submit", async event => {
    event.preventDefault();
    if (authMode === "login") {
        await loginUser();
    } else {
        await register();
    }
});

/* 
LOGOUT
 */

function logout() {
    clearCredentials();
    currentGame = null;
    currentGameId = null;
    board.innerHTML = "";
    loginInput.value = "";
    passwordInput.value = "";
    setAuthStatus("");
    showAuth();
}

logoutBtn.addEventListener("click", logout);

/* 
HANDLE 401
 */

function handleUnauthorized() {
    clearCredentials();
    currentGame = null;
    currentGameId = null;
    showAuth();
    setAuthStatus("Сессия авторизации недействительна.", "error");
}

/* 
CREATE COMPUTER GAME
 */

async function createComputerGame() {
    try {
        const response = await fetchWithAuth(`${API_BASE}/game/new?type=computer`, {
            method: "GET"
        });

        if (response.status === 401) return;

        const data = await parseResponse(response);

        if (!response.ok) {
            throw new Error(getErrorMessage(data, "Не удалось создать игру"));
        }

        currentGameId = data.id || data.uuid || data.UUID;
        await openGame(currentGameId);

    } catch (error) {
        console.error("Create computer game error:", error);
        setMenuStatus(error.message, "error");
    }
}

/* 
CREATE PLAYER GAME
 */

async function createPlayerGame() {
    try {
        const response = await fetchWithAuth(`${API_BASE}/game/new?type=player`, {
            method: "GET"
        });

        if (response.status === 401) return;

        const data = await parseResponse(response);

        if (!response.ok) {
            throw new Error(getErrorMessage(data, "Не удалось создать игру"));
        }

        currentGameId = data.id || data.uuid || data.UUID;
        await openGame(currentGameId);

    } catch (error) {
        console.error("Create player game error:", error);
        setMenuStatus(error.message, "error");
    }
}

/* 
AVAILABLE GAMES
 */

async function loadAvailableGames() {
    try {
        hideExtraSection();
        gamesList.innerHTML = `<p class="empty">Загрузка...</p>`;

        const response = await fetchWithAuth(`${API_BASE}/game/available`, {
            method: "GET"
        });

        if (response.status === 401) return;

        const data = await parseResponse(response);

        if (!response.ok) {
            throw new Error(getErrorMessage(data, "Не удалось получить список игр"));
        }

        renderAvailableGames(Array.isArray(data) ? data : []);

    } catch (error) {
        console.error("Load games error:", error);
        gamesList.innerHTML = `<p class="empty">Не удалось загрузить игры</p>`;
        setMenuStatus(error.message, "error");
    }
}

function renderAvailableGames(games) {
    gamesList.innerHTML = "";

    if (games.length === 0) {
        gamesList.innerHTML = `<p class="empty">Нет доступных игр</p>`;
        return;
    }

    games.forEach(game => {
        const gameId = game.uuid || game.id || game.UUID;

        const item = document.createElement("div");
        item.className = "game-item";

        const info = document.createElement("div");
        info.className = "game-item-info";

        const title = document.createElement("div");
        title.className = "game-item-title";
        title.textContent = "Игра с игроком";

        const id = document.createElement("div");
        id.className = "game-item-id";
        id.textContent = gameId || "Неизвестная игра";

        info.appendChild(title);
        info.appendChild(id);

        const joinButton = document.createElement("button");
        joinButton.textContent = "Присоединиться";

        joinButton.addEventListener("click", () => joinGame(gameId));

        item.appendChild(info);
        item.appendChild(joinButton);

        gamesList.appendChild(item);
    });
}

async function showGameHistory() {
    try {
        hideExtraSection();

        const response = await fetchWithAuth(`${API_BASE}/game/history`, {
            method: "GET"
        });

        if (response.status === 401) return;

        const data = await parseResponse(response);

        if (!response.ok) {
            throw new Error(getErrorMessage(data, "Не удалось загрузить историю игр"));
        }

        const games = Array.isArray(data) ? data : [];

        if (games.length === 0) {
            showExtraSection("История игр", `<p class="empty">История игр пуста</p>`);
            return;
        }

        const rows = games.map(game => {
            const status = game.status || game.Status;
            const createdAt = game.created_at || game.CreatedAt || "";
            return `
                <div class="extra-item">
                    <div><strong>ID:</strong> ${game.uuid || game.id || game.UUID || "-"}</div>
                    <div><strong>Статус:</strong> ${status}</div>
                    <div><strong>Дата:</strong> ${createdAt}</div>
                </div>`;
        }).join("");

        showExtraSection("История игр", `<div class="extra-list">${rows}</div>`);
    } catch (error) {
        console.error("Load history error:", error);
        showExtraSection("История игр", `<p class="empty">${error.message}</p>`);
    }
}

async function showLeaderboard() {
    try {
        hideExtraSection();

        const response = await fetchWithAuth(`${API_BASE}/user/leaderboard`, {
            method: "GET"
        });

        if (response.status === 401) return;

        const data = await parseResponse(response);

        if (!response.ok) {
            throw new Error(getErrorMessage(data, "Не удалось загрузить таблицу лидеров"));
        }

        const entries = Array.isArray(data) ? data : [];

        if (entries.length === 0) {
            showExtraSection("Таблица лидеров", `<p class="empty">Лидерборд пуст</p>`);
            return;
        }

        const rows = entries.map(entry => {
            return `
                <div class="extra-item">
                    <div><strong>${entry.login || entry.Login}</strong></div>
                    <div>Рейтинг: ${entry.win_ratio || entry.WinRatio || 0}</div>
                </div>`;
        }).join("");

        showExtraSection("Таблица лидеров", `<div class="extra-list">${rows}</div>`);
    } catch (error) {
        console.error("Load leaderboard error:", error);
        showExtraSection("Таблица лидеров", `<p class="empty">${error.message}</p>`);
    }
}

/* 
JOIN GAME
 */

async function joinGame(id) {
    if (!id) {
        setMenuStatus("Некорректный ID игры", "error");
        return;
    }

    try {
        setMenuStatus("Подключаемся к игре...", "info");

        const response = await fetchWithAuth(`${API_BASE}/game/join/${id}`, {
            method: "POST"
        });

        if (response.status === 401) return;

        const data = await parseResponse(response);

        if (!response.ok) {
            throw new Error(getErrorMessage(data, "Не удалось присоединиться к игре"));
        }

        currentGameId = id;
        await openGame(currentGameId);

    } catch (error) {
        console.error("Join game error:", error);
        setMenuStatus(error.message, "error");
    }
}

/* 
OPEN GAME
 */

async function openGame(id) {
    currentGameId = id;
    showGame();
    gameIdElement.textContent = `ID: ${id}`;
    await loadGame();
}

/* 
LOAD GAME
 */

async function loadGame() {
    if (!currentGameId) return;

    try {
        gameStatus.textContent = "Загрузка...";

        const response = await fetchWithAuth(`${API_BASE}/game/info/${currentGameId}`, {
            method: "GET"
        });

        if (response.status === 401) return;

        const data = await parseResponse(response);

        if (!response.ok) {
            throw new Error(getErrorMessage(data, "Игра не найдена"));
        }

        currentGame = data;
        renderBoard();
        updateGameStatus();

    } catch (error) {
        console.error("Load game error:", error);
        gameStatus.textContent = error.message;
    }
}

/* 
BOARD
 */

function getBoard() {
    if (currentGame) {
        return currentGame.playing_field || currentGame.board || currentGame.PlayingField;
    }
    return null;
}

function getSymbol(value) {
    if (value === 1) return "X";
    if (value === 2) return "O";
    return "";
}

function renderBoard() {
    board.innerHTML = "";
    const field = getBoard();

    if (!field) {
        gameStatus.textContent = "Поле игры недоступно";
        return;
    }

    for (let row = 0; row < 3; row++) {
        for (let column = 0; column < 3; column++) {
            const cell = document.createElement("div");
            const value = field[row][column];
            const symbol = getSymbol(value);

            cell.className = "cell";

            if (symbol === "X") cell.classList.add("x");
            if (symbol === "O") cell.classList.add("o");

            if (!symbol) {
                cell.classList.add("empty");
                if (canMakeMove()) {
                    cell.addEventListener("click", () => makeMove(row, column));
                }
            }

            cell.textContent = symbol;
            board.appendChild(cell);
        }
    }
}

/* 
CAN MAKE MOVE
 */

function canMakeMove() {
    if (!currentGame) return false;

    const status = currentGame.status || currentGame.Status;

    if (status === "X_WON" || status === "O_WON" || status === "DRAW") {
        return false;
    }

    const gameType = currentGame.game_type || currentGame.GameType;

    if (gameType === "COMPUTER") {
        const userX = currentGame.user_x || currentGame.UserX;
        const currentPlayer = currentGame.current_player || currentGame.CurrentPlayer;
        return currentPlayer === userX;
    }

    const userUUID = localStorage.getItem("userUUID");
    const currentPlayer = currentGame.current_player || currentGame.CurrentPlayer;

    if (!userUUID) {
        return status === "IN_PROGRESS" || status === "X_TURN" || status === "O_TURN";
    }

    return currentPlayer === userUUID;
}

/* 
MAKE MOVE
 */

async function makeMove(row, column) {
    if (!currentGameId) return;

    try {
        const field = getBoard();

        if (!field || field[row][column] !== 0) {
            return;
        }

        gameStatus.textContent = "Выполняется ход...";

        const response = await fetchWithAuth(`${API_BASE}/game/move/${currentGameId}`, {
            method: "POST",
            body: JSON.stringify({ row, column })
        });

        if (response.status === 401) return;

        const data = await parseResponse(response);

        if (!response.ok) {
            throw new Error(getErrorMessage(data, "Ошибка выполнения хода"));
        }

        currentGame = data;
        renderBoard();
        updateGameStatus();

    } catch (error) {
        console.error("Move error:", error);
        gameStatus.textContent = error.message;
        await loadGame();
    }
}

/* 
GAME STATUS
 */

function updateGameStatus() {
    if (!currentGame) return;

    const status = currentGame.status || currentGame.Status;
    gameStatus.style.color = "";

    switch (status) {
        case "IN_PROGRESS":
        case "X_TURN":
            gameStatus.textContent = "Ход X";
            gameStatus.style.color = "#4f46e5";
            break;

        case "O_TURN":
            gameStatus.textContent = "Ход O";
            gameStatus.style.color = "#dc2626";
            break;

        case "X_WON":
            gameStatus.textContent = "Победил X!";
            gameStatus.style.color = "#16a34a";
            break;

        case "O_WON":
            gameStatus.textContent = "Победил O!";
            gameStatus.style.color = "#16a34a";
            break;

        case "DRAW":
            gameStatus.textContent = "Ничья!";
            gameStatus.style.color = "#d97706";
            break;

        default:
            gameStatus.textContent = `Статус: ${status}`;
    }
}

/* 
MENU BUTTONS
 */

computerGameBtn.addEventListener("click", createComputerGame);
playerGameBtn.addEventListener("click", createPlayerGame);
refreshGamesBtn.addEventListener("click", loadAvailableGames);

/* 
GAME BUTTONS
 */

refreshGameBtn.addEventListener("click", loadGame);
historyBtn.addEventListener("click", showGameHistory);
leaderboardBtn.addEventListener("click", showLeaderboard);

backBtn.addEventListener("click", () => {
    currentGame = null;
    currentGameId = null;
    board.innerHTML = "";
    showMenu();
});

/* 
INIT
 */

async function init() {
    const tokens = getTokens();

    if (!tokens) {
        showAuth();
        return;
    }

    try {
        const response = await fetchWithAuth(`${API_BASE}/user/CurrentUser`, {
            method: "GET"
        });

        if (response.status === 401 || !response.ok) {
            clearCredentials();
            showAuth();
            return;
        }

        const data = await parseResponse(response);

        if (data && typeof data === "object") {
            if (data.login) {
                localStorage.setItem("userLogin", data.login);
            }

            const uuid = data.uuid || data.UUID || data.id;
            if (uuid) {
                localStorage.setItem("userUUID", uuid);
            }
        }

        showMenu();

    } catch (error) {
        console.error("Initialization error:", error);
        showAuth();
        setAuthStatus("Не удалось подключиться к серверу", "error");
    }
}

init();