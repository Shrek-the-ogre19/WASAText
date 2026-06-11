import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5
});

instance.interceptors.request.use((config) => {
	const token = localStorage.getItem("identifier");
	if (token) {
		config.headers.Authorization = `Bearer ${token}`;
	}
	return config;
});

export const REFRESH_INTERVAL_MS = 3000;

export function startAutoRefresh(callback) {
	const timer = setInterval(callback, REFRESH_INTERVAL_MS);
	return () => clearInterval(timer);
}

export const GROUP_DEFAULT_PICTURE = "/9572728.png";
export const USER_DEFAULT_PICTURE = "/default-avatar-icon-of-social-media-user-vector.jpg";

export function isDefaultPicture(picture) {
	if (picture == null || picture === "" || picture === "default") {
		return true;
	}
	return (
		!picture.startsWith("data:") &&
		!picture.startsWith("http://") &&
		!picture.startsWith("https://") &&
		!picture.startsWith("/")
	);
}

export default instance;
