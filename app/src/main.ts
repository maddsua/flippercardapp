import { createApp } from 'vue';
import { createWebHistory, createRouter } from 'vue-router';

import App from './App.vue';
import PlayView from './components/Play/PlayView.vue';
import HomeView from './components/Home/HomeView.vue';
import CollectionsListView from './components/Collections/CollectionsListView.vue';
import CollectionView from './components/Collections/CollectionView.vue';

import './main.scss';
import './theme.scss';
import DiscoverView from './components/Discover/DiscoverView.vue';
import StarredView from './components/Starred/StarredView.vue';
import MenuView from './components/Menu/MenuView.vue';

const routes = [
	{
		path: '/',
		component: HomeView,
	},
	{
		path: '/app',
		component: CollectionsListView,
		meta: {
			app_view: 'home'
		},
	},
	{
		path: '/app/collections',
		component: CollectionsListView,
		meta: {
			app_view: 'home'
		},
	},
	{
		path: '/app/collection/:collection_id',
		component: CollectionView,
		meta: {
			app_view: 'home'
		},
	},
	{
		path: '/app/play/deck/:deck_id',
		component: PlayView,
	},
	{
		path: '/app/discover',
		component: DiscoverView,
	},
	{
		path: '/app/starred',
		component: StarredView,
	},
	{
		path: '/app/menu',
		component: MenuView,
	},
	{
		path: '/:pathMatch(.*)*',
		component: HomeView,
	},
]

const router = createRouter({
	history: createWebHistory(),
	routes,
});

createApp(App).use(router).mount('#app-root')
