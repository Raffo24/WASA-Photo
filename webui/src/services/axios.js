import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 50
});

const axiosUpdate = () => {
	instance.defaults.headers.common['Authorization'] = 'Bearer ' + getCurrentSession();
}
function getCurrentSession() {
	if (localStorage.getItem('token') == null){
        return sessionStorage.getItem('token');
    }
    return localStorage.getItem('token');
}
export {
	instance as axios,
	axiosUpdate as axiosUpdate,
	getCurrentSession as getCurrentSession
}