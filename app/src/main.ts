import { createApp } from 'vue';
import { createRouter, createWebHistory } from 'vue-router';

import './main.scss';
import './theme.scss';

import { useClient } from './api';
import App from './App.vue';
import NotFoundView from './components/App/NotFoundView.vue';
import AllCollectionsView from './components/Collections/AllCollections/AllCollectionsView.vue';
import EditCollectionView from './components/Collections/CollectionEditor/EditCollectionView.vue';
import NewCollectionView from './components/Collections/CollectionEditor/NewCollectionView.vue';
import CollectionView from './components/Collections/CollectionView.vue';
import DiscoverView from './components/Collections/DiscoverCollections/DiscoverView.vue';
import MyCollectionsView from './components/Collections/MyCollectionsView.vue';
import DashboardSettingsScreen from './components/Dashboard/Settings/DashboardSettingsScreen.vue';
import DashboardView from './components/Dashboard/DashboardView.vue';
import DeckEditorView from './components/Decks/Editor/DeckEditorView.vue';
import PlayView from './components/Play/PlayView.vue';
import StarredView from './components/Starred/StarredView.vue';
import DashboardSigninScreen from './components/Dashboard/Auth/DashboardSigninScreen.vue';

const client = useClient();

const routes = [
	{
		path: '/',
		component: MyCollectionsView,
		meta: {
			app_view: 'home'
		},
	},
	{
		path: '/collections/discover',
		component: DiscoverView,
		meta: {
			app_view: 'discover'
		},
	},
	{
		path: '/collections',
		redirect: () => {
			return { path: '/collections/discover' };
		},
	},
	{
		path: '/discover',
		redirect: () => {
			return { path: '/collections/discover' };
		},
	},
	{
		path: '/collections/all',
		component: AllCollectionsView,
		meta: {
			app_view: 'discover'
		},
	},
	{
		path: '/collections/new',
		component: NewCollectionView,
		meta: {
			app_view: 'menu',
			requiresDashboardSession: true,
			requiresEditorPermission: true,
		},
	},
	{
		path: '/collection/:collection_id',
		component: CollectionView,
		meta: {
			app_view: 'home'
		},
	},
	{
		path: '/collection/:collection_id/edit',
		component: EditCollectionView,
		meta: {
			app_view: 'menu',

		},
	},
	{
		path: '/play/deck/:deck_id',
		component: PlayView,
	},
	{
		path: '/starred',
		component: StarredView,
		meta: {
			app_view: 'starred'
		},
	},
	{
		path: `/decks/editor`,
		component: DeckEditorView,
		meta: {
			requiresDashboardSession: true,
			requiresEditorPermission: true,
		},
	},
	{
		path: `/decks/editor/:deck_id`,
		component: DeckEditorView,
		meta: {
			requiresDashboardSession: true,
			requiresEditorPermission: true,
		},
	},
	{
		path: '/dashboard',
		component: DashboardView,
		meta: {
			app_view: 'menu'
		},
		children: [
			{
				path: '',
				component: DashboardSettingsScreen,
			},
			{
				path: 'settings',
				component: DashboardSettingsScreen,
			},
			{
				path: 'signin',
				component: DashboardSigninScreen,
			}
		],
	},
	{
		path: '/:pathMatch(.*)*',
		component: NotFoundView,
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
		return '/dashboard/signin';
	} else if (requiresEditorPermission && !authState?.actor?.permissions.content_edit) {
		console.warn('ROUTER: Content editor permission missing. Redirecting to the dashboard home');
		return '/dashboard';
	}

});

createApp(App).use(router).mount('#app-root')
