<script>
export default {
    data: function () {
        return {
            errormsg: null,
            loading: false,
            field_username: "",
            rememberLogin: false,
        };
    },
    methods: {
        async login() {
            this.loading = true;
            this.errormsg = null;
        
            let response = await this.$axios.post("/session", 
                {
                    username: this.field_username
                }   
            );

            if (response == null) {
				this.loading = false
				return
			}

            // If the login is successful, save the token and redirect to the previous page
            if (response.status == 201 || response.status == 200) {
                // Save the token in the local storage if the user wants to be remembered
                if (this.rememberLogin) {
                    localStorage.setItem("token", response.data["ID"])
                    sessionStorage.removeItem("token");
                }
                // Else save the token in the session storage
                else {
                    sessionStorage.setItem("token", response.data["ID"]);
                    localStorage.removeItem("token");
                }
                // Tell the root view to enable the navbar
                this.$root.setLoggedIn();
                // Update the header
                this.$axiosUpdate();

                // Go back to the previous page
                this.$router.go(-1);
            }
            else {
                // Login failed, show the error message
                this.errormsg = response.data["error"];
            }
            this.loading = false;
        },
    },
}
</script>

<template>
    <div class="vh-100 container py-5 h-100">
        <div class="row d-flex justify-content-center align-items-center h-100">
            <div class="col-12 col-md-8 col-lg-6 col-xl-5">
                <div class="card bg-dark text-white" style="border-radius: 1rem">
                    <div class="card-body p-4">

                        <h1 class="h2 pb-4 text-center">WASAPhoto</h1>

                        <form>
                            <div class="form-floating mb-4">
                                <input v-model="field_username" type="email" id="formUsername" class="form-control bg-dark text-white"
                                    placeholder="name@example.com" />
                                <label class="form-label" for="formUsername">Username</label>
                            </div>

                            <!-- 2 column grid layout for inline styling -->
                            <div class="row mb-4">
                                <div class="col d-flex justify-content-center">
                                    <!-- Checkbox -->
                                    <div class="form-check">
                                        <input v-model="rememberLogin" class="form-check-input" type="checkbox" value=""
                                            id="form2" />
                                        <label class="form-check-label" for="form2">Remember me</label>
                                    </div>
                                </div>
                            </div>

                            <!-- Submit button -->
                            <button style="width: 100%" type="button" class="btn btn-outline-primary btn-block mb-4"
                                @click="login">Sign in</button>
                            <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
                            <LoadingSpinner :loading="loading" />
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style>

</style>
