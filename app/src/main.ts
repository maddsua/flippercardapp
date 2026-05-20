import { createApp } from 'vue';
import { createWebHistory, createRouter } from 'vue-router';

import './main.scss';
import './theme.scss';

import App from './App.vue';
import PlayView from './components/Play/PlayView.vue';
import HomeView from './components/Home/HomeView.vue';
import CollectionsListView from './components/Collections/CollectionsListView.vue';
import CollectionView from './components/Collections/CollectionView.vue';
import DiscoverView from './components/Discover/DiscoverView.vue';
import StarredView from './components/Starred/StarredView.vue';
import DashboardView from './components/Dashboard/DashboardView.vue';
import DashboardMainScreen from './components/Dashboard/DashboardMainScreen.vue';
import DashboardContentScreen from './components/Dashboard/Content/DashboardContentScreen.vue';
import NewCollectionScreen from './components/Dashboard/Content/Collections/NewCollectionScreen.vue';
import EditCollectionMetadataScreen from './components/Dashboard/Content/Collections/EditCollectionMetadataScreen.vue';
import DeckEditorView from './components/DeckEditor/DeckEditorView.vue';
import DashboardCollectionScreen from './components/Dashboard/Content/Collections/DashboardCollectionScreen.vue';
import { useClient } from './api';

const client = useClient();

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
			{
				path: '',
				component: DashboardMainScreen,
			},
			{
				path: 'content',
				component: DashboardContentScreen,
				meta: {
					requiresDashboardSession: true,
					requiresEditorPermission: true
				},
			},
			{
				path: 'content/collections/new',
				component: NewCollectionScreen,
				meta: {
					requiresDashboardSession: true,
					requiresEditorPermission: true
				},
			},
			{
				path: 'content/collection/:collection_id/metadata',
				component: EditCollectionMetadataScreen,
				meta: {
					requiresDashboardSession: true,
					requiresEditorPermission: true
				},
			},
			{
				path: 'content/collection/:collection_id',
				component: DashboardCollectionScreen,
				meta: {
					requiresDashboardSession: true,
					requiresEditorPermission: true
				},
			},
		],
	},
	{
		path: `/app/editor/deck/:deck_id/editor`,
		component: DeckEditorView,
		meta: {
			requiresDashboardSession: true,
			requiresEditorPermission: true
		},
	},
	{
		path: `/app/editor/deck/editor`,
		component: DeckEditorView,
		meta: {
			requiresDashboardSession: true,
			requiresEditorPermission: true
		},
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

router.beforeResolve(async (to) => {

	const { requiresDashboardSession, requiresEditorPermission } = to.meta;

	const { data: authState } = await client.auth.whoami({ cached: true });

	if (requiresDashboardSession && !authState?.actor) {
		console.warn('ROUTER: Unauthorized. Redirecting to the login screen');
		return '/app/dashboard';
	} else if (requiresEditorPermission && !authState?.actor?.permissions.content_edit) {
		console.warn('ROUTER: Content editor permission missing. Redirecting to the dashboard home');
		return '/app/dashboard';
	}

});

createApp(App).use(router).mount('#app-root')
