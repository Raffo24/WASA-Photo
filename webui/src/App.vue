<script>
export default {
	props: ["user_id", "name", "date", "comments", "likes", "photo_id", "liked"],
	data: function () {
		return {
			modalTitle: "Modal Title",
			modalMsg: "Modal Message",

			logged_in: true,
		}
	},
	methods: {

		showModal(title, message) {
			this.modalTitle = title
			this.modalMsg = message

			this.$refs.errModal.showModal()
		},

		setLoggedIn() {
			this.logged_in = true
		},

		logout() {
			localStorage.removeItem("token")
            sessionStorage.removeItem("token")
			this.logged_in = false
            this.$router.push({ path: "/login" })
		}
	},

	mounted() {
		this.$axiosUpdate()

		this.$axios.interceptors.response.use(response => {
			return response;
		}, error => {
			if (error.response.status != 0) {
				if (error.response.status === 401) {
					this.$router.push({ path: '/login' })
					this.logged_in = false;
					return;
				}
				
				this.showModal("Error " + error.response.status, error.response.data['status'])
				return;
			}
			this.showModal("Error", error.toString());
			return;
		});
	}
}
</script>

<template>
	<body style="background-color: #161819">
		<Modal ref="errModal" id="errorModal" :title="modalTitle">
			{{ modalMsg }}
		</Modal>
		<div class="container-fluid" id="app">
			<div class="row">
				<main>
					<RouterView />
					<div v-if="logged_in" class="mb-5 pb-3"></div>
				</main>

				<!-- Navigation bar -->
				<nav v-if="logged_in" id="global-nav" class="navbar fixed-bottom navbar-expand-lg navbar-dark bg-dark"> 
					<div class="collapse navbar-collapse" id="navbarNav"></div>
					<RouterLink to="/" class="col-4 text-center">
						<i class="bi bi-house text-light" style="font-size: 2em"></i>
					</RouterLink>
					<RouterLink to="/search" class="col-4 text-center">
						<i class="bi bi-search text-light" style="font-size: 2em"></i>
					</RouterLink>
					<RouterLink to="/users/me" class="col-4 text-center">
						<i class="bi bi-person text-light" style="font-size: 2em"></i>
					</RouterLink>
				</nav>
			</div>
		</div>
	</body>
</template>
