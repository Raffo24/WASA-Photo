<script>
import ProfileCounters from "./ProfileCounters.vue"
export default {
    props: ["user_id", "name", "followed", "banned", "show_new_post", "user_data"],
    components: {
        ProfileCounters
    },
    watch: {
        name: function (new_val, old_val) { this.username = new_val },
        banned: function (new_val, old_val) { this.user_banned = new_val },
        followed: function (new_val, old_val) { this.user_followed = new_val },
        user_id: function (new_val, old_val) { this.myself = this.$currentSession() == new_val },
    },
    data: function () {
        return {
            // User data
            username: this.name,
            user_followed: this.followed,
            user_banned: this.banned,
            formData: new FormData(),
            user_datas: this.user_data,

            myself: this.$currentSession() == this.user_id,

            show_post_form: false,
            show_username_form: false,

            newUsername: "",

            upload_file: null,
            title: "",
            description: "",
        }
    },
    methods: {
        logout() {
            this.$root.logout()
        },

        visit() {
            this.$router.push({ path: "/users/" + this.user_id });
        },

        follow() {
            this.$axios.put("/users/" + this.$currentSession()+ "/follow/" + this.user_id)
                .then(response => {
                    if (response == null) return
                    this.user_followed = true
                    this.$emit('updateInfo')
                })
        },

        unfollow() {
            this.$axios.delete("/users/" + this.$currentSession()+ "/follow/" + this.user_id)
                .then(response => {
                    if (response == null) return
                    this.user_followed = false
                    this.$emit('updateInfo')
                })
        },

        ban() {
            this.$axios.put("/users/" + this.$currentSession() + "/ban/" + this.user_id)
                .then(response => {
                    if (response == null) return
                    this.user_banned = true
                    this.$emit('updateInfo')
                })
        },

        unban() {
            this.$axios.delete("/users/" + this.$currentSession() + "/ban/" + this.user_id)
                .then(response => {
                    if (response == null) return
                    this.user_banned = false
                    this.$emit('updateInfo')
                })
        },

        load_file(e) {
            let files = e.target.files || e.dataTransfer.files;
            if (!files.length) return
            this.upload_file = files[0]
        },

        submit_file() {
            let formData = new FormData();
            formData.append('photo', this.upload_file, this.upload_file.name);          
            formData.append('title', this.title);
            formData.append('description', this.description);
            // send the request
            this.$axios.post("/photos", formData)
                .then(response => {
                    if (response == null) return
                    this.show_post_form = false
                    this.$emit('updatePosts')
                })
        },

        updateUsername() {
            this.$axios.put("/users/" + this.$currentSession(), { Username: this.newUsername })
                .then(response => {
                    if (response == null) return
                    this.show_username_form = false
                    this.$emit('updateInfo')
                    this.username = this.newUsername
                })
        },
    },
}
</script>
<template>
    
    <div class="card mb-4 bg-dark text-white ">
        <div class="container">
            <div class="row">
                <div class="col-12">
                    <div class="card-body h-100" :style="{
                                            'justify-content': (user_datas ? 'center' : 'space-between'),
                                            'align-items': (user_datas ? 'center' : 'start'),
                                            'display': 'flex'
                                            }">
                        <a @click="visit">
                            <h5 class="card-title mb-0 d-inline-block" style="cursor: pointer">
                                {{ username }} 
                            </h5>
                        </a>
                        
                        <div class="d-flex flex-column" v-if="!(user_datas)">
                            <div v-if="!myself" class="d-flex">
                                <button v-if="!user_banned" @click="ban" type="button"
                                    class="btn btn-outline-danger me-2">Ban</button>
                                <button v-if="user_banned" @click="unban" type="button"
                                    class="btn btn-outline-danger me-2">Banned</button>
                                <button v-if="!user_followed" @click="follow" type="button"
                                    class="btn btn-outline-primary">Follow</button>
                                <button v-if="user_followed" @click="unfollow" type="button"
                                    class="btn btn-outline-primary">Following</button>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="col-12" v-if="(user_datas)">
                    <ProfileCounters :user_data="user_datas" />
                </div>
                <div class="d-flex flex-column" v-if="(user_datas)" v-bind:class="{
                                                'col-12': (myself && show_new_post),
                                                'col-sm:12': (myself && show_new_post),
                                                'col-12': !(myself && show_new_post),
                                                'align-items-center': !(myself && show_new_post),
                                                'align-items-sm-center': (myself && show_new_post),
                                            }">
                    <!-- Buttons -->
                    <div class="card-body d-flex" >
                        <div v-if="!myself" class="d-flex">
                            <button v-if="!user_banned" @click="ban" type="button"
                                class="btn btn-outline-danger me-2">Ban</button>
                            <button v-if="user_banned" @click="unban" type="button"
                                class="btn btn-outline-danger me-2">Banned</button>
                            <button v-if="!user_followed" @click="follow" type="button"
                                class="btn btn-outline-primary">Follow</button>
                            <button v-if="user_followed" @click="unfollow" type="button"
                                class="btn btn-outline-primary">Following</button>
                        </div>
                        
                        <!-- Users cannot follow or ban themselves -->
                        <div v-if="(myself && !show_new_post)" >
                            <button disabled type="button" class="btn btn-secondary">Yourself</button>
                        </div>

                        <div class="d-flex col justify-content-center flex-row">
                            <!-- Logout button -->
                            <div v-if="(myself && show_new_post)">
                                <button type="button" class="btn btn-outline-danger me-2" @click="logout">Logout</button>
                            </div>
                            <!-- Update username button -->
                            <div v-if="(myself && show_new_post)">
                                <button v-if="!show_username_form" type="button" class="btn btn-outline-secondary me-2"
                                    @click="show_username_form = true">Username</button>
                            </div>

                            <!-- Post a new photo button -->
                            <div v-if="(myself && show_new_post)">
                                <button v-if="!show_post_form" type="button" class="btn btn-outline-primary me-2"
                                    @click="show_post_form = true">Post</button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- File input -->
            <div class="row bg-dark text-white" v-if="show_post_form">
                <!-- title input-->
                <div  class="col-4">
                    <input v-model="title" class="form-control form-control-lg bg-dark text-white" id="formTitle" placeholder="Titolo" />
                    <!-- description input-->
                </div>
                <div class="col-8">
                    <input v-model="description" class="form-control form-control-lg bg-dark text-white" id="formDescription"
                        placeholder="Descrizione" />
                </div>
                <div class="col-9">
                    <div class="card-body h-100 d-flex align-items-center">
                        <input @change="load_file" class="form-control form-control-lg bg-dark text-white" id="formFileLg" type="file" />
                    </div>
                </div>

                <!-- Publish button -->
                <div class="col-3">
                    <div class="card-body d-flex justify-content-center">
                        <button type="button" class="btn btn-outline-primary btn-lg" @click="submit_file">Publish</button>
                    </div>
                </div>
            </div>

            <!-- Change username form -->
            <div class="row bg-dark text-white" v-if="show_username_form">

                <!-- Username input -->
                <div class="col-10">
                    <div class="card-body h-100 d-flex align-items-center">
                        <input v-model="newUsername" class="form-control form-control-lg bg-dark text-white" id="formUsername"
                            placeholder="Il tuo nuovo username" />
                    </div>
                </div>

                <!-- Username update button -->
                <div class="col-2">
                    <div class="card-body d-flex justify-content-center">
                        <button type="button" class="btn btn-outline-primary btn-lg" @click="updateUsername">Set</button>
                    </div>
                </div>
            </div>
        </div>
    </div>

</template>