import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from "../views/LoginView.vue";
import MainpageView from "../views/MainpageView.vue";
import ProfileView from "../views/ProfileView.vue";
import UsersListView from "../views/UsersListView.vue";
import ConversationView from "../views/ConversationView.vue";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: HomeView},
		{path: '/link1', component: HomeView},
		{path: '/link2', component: HomeView},
		{path: '/some/:id/link', component: HomeView},


		{path: "/session", component: LoginView},
		{path: "/mainpage/:Id/conversations", component: MainpageView, props: true},
		{path: "/mainpage/:Id", component: ProfileView},
		{path: "/:Id/users", component: UsersListView},
		{path: "/:Id/conversations/:conversationId", component: ConversationView}

	]
})

export default router
