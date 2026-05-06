import { createApp } from 'vue';
import { createWebHistory, createRouter } from 'vue-router';

import App from './App.vue';
import PlayView from './components/Play/PlayView.vue';
import HomeView from './components/Home/HomeView.vue';
import CollectionsListView from './components/Collections/CollectionsListView.vue';
import CollectionView from './components/Collections/CollectionView.vue';

import './main.scss';
import './theme.scss';

const routes = [
	{ path: '/', component: HomeView },
	{ path: '/app', component: CollectionsListView },
	{ path: '/app/collections', component: CollectionsListView },
	{ path: '/app/collection/:collection_id', component: CollectionView },
	{ path: '/app/play/deck/:deck_id', component: PlayView },
	{ path: '/:pathMatch(.*)*', component: HomeView },
]

const router = createRouter({
	history: createWebHistory(),
	routes,
});

createApp(App).use(router).mount('#app-root')
