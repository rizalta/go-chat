export function updateTimeAgo(givenTime) {
	const now = new Date();
	const diffInSeconds = Math.floor((now - givenTime) / 1000);

	if (diffInSeconds < 60) {
		return `${diffInSeconds} seconds ago`;
	} else if (diffInSeconds < 3600) {
		return `${Math.floor(diffInSeconds / 60)} minutes ago`;
	} else if (diffInSeconds < 86400) {
		return `${Math.floor(diffInSeconds / 3600)} hours ago`;
	} else {
		return `${Math.floor(diffInSeconds / 86400)} days ago`;
	}
}
