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
import DashboardView from './components/Dashboard/DashboardView.vue';
import DashboardMainScreen from './components/Dashboard/DashboardMainScreen.vue';
import DashboardContentScreen from './components/Dashboard/Content/DashboardContentScreen.vue';
import NewCollectionScreen from './components/Dashboard/Content/Collections/NewCollectionScreen.vue';
import EditCollectionMetadataScreen from './components/Dashboard/Content/Collections/EditCollectionMetadataScreen.vue';
import DeckEditorView from './components/DeckEditor/DeckEditorView.vue';
import DashboardCollectionScreen from './components/Dashboard/Content/Collections/DashboardCollectionScreen.vue';

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
		meta: {
			app_view: 'discover'
		},
	},
	{
		path: '/app/starred',
		component: StarredView,
		meta: {
			app_view: 'starred'
		},
	},
	{
		path: '/app/dashboard',
		component: DashboardView,
		meta: {
			app_view: 'menu'
		},
		children: [
			{ path: '', component: DashboardMainScreen },
			{ path: 'content', component: DashboardContentScreen },
			{ path: 'content/collections/new', component: NewCollectionScreen },
			{ path: 'content/collection/:collection_id/metadata', component: EditCollectionMetadataScreen },
			{ path: 'content/collection/:collection_id', component: DashboardCollectionScreen },
		],
	},
	{
		path: `/app/editor/deck/:deck_id/editor`,
		component: DeckEditorView,
	},
	{
		path: `/app/editor/deck/editor`,
		component: DeckEditorView,
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
