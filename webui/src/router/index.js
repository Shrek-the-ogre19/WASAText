import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from "../views/LoginView.vue";
import MainpageView from "../views/MainpageView.vue";
import ProfileView from "../views/ProfileView.vue";
import UsersListView from "../views/UsersListView.vue";
import ConversationView from "../views/ConversationView.vue";
import About from "../views/About.vue";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		/*{path: '/', component: HomeView},
		{path: '/link1', component: HomeView},
		{path: '/link2', component: HomeView},
		{path: '/some/:id/link', component: HomeView},

*/

		{path: "", component: About},
		{path: "/session", component: LoginView},
		{path: "/users", component: UsersListView},
		{path: "/mainpage/:Id/conversations", component: MainpageView, props: true},
		{path: "/mainpage/:Id", component: ProfileView},
		{path: "/mainpage/:Id/conversations/:conversationId", component: ConversationView}

	]
})

export default router
