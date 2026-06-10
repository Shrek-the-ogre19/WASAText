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
