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
import { appSetTitle, enablePwaInstall } from './app';

const client = useClient();

const routes = [
	{
		path: '/',
		component: MyCollectionsView,
		meta: {
			app_view: 'home',
			app_title: 'My cards',
		},
	},
	{
		path: '/collections/discover',
		component: DiscoverView,
		meta: {
			app_view: 'discover',
			app_title: 'Discover',
		},
	},
	{
		path: '/collections',
		redirect: () => ({ path: '/collections/discover' }),
	},
	{
		path: '/discover',
		redirect: () => ({ path: '/collections/discover' }),
	},
	{
		path: '/collections/all',
		component: AllCollectionsView,
		meta: {
			app_view: 'discover',
			app_title: 'All collections',
		},
	},
	{
		path: '/collections/new',
		component: NewCollectionView,
		meta: {
			app_view: 'menu',
			app_title: 'Create collection',
			requiresDashboardSession: true,
			requiresEditorPermission: true,
		},
	},
	{
		path: '/collection/:collection_id',
		component: CollectionView,
		meta: {
			app_view: 'home',
			app_title: 'Card collection',
		},
	},
	{
		path: '/collection/:collection_id/edit',
		component: EditCollectionView,
		meta: {
			app_view: 'menu',
			app_title: 'Edit card collection',
		},
	},
	{
		path: '/play/deck/:deck_id',
		component: PlayView,
		app_title: 'Play deck',
	},
	{
		path: '/starred',
		component: StarredView,
		meta: {
			app_view: 'starred',
			app_title: 'Starred decks',
		},
	},
	{
		path: `/decks/editor`,
		component: DeckEditorView,
		meta: {
			app_title: 'Deck editor',
			requiresDashboardSession: true,
			requiresEditorPermission: true,
		},
	},
	{
		path: `/decks/editor/:deck_id`,
		component: DeckEditorView,
		meta: {
			app_title: 'Deck editor',
			requiresDashboardSession: true,
			requiresEditorPermission: true,
		},
	},
	{
		path: '/dashboard',
		component: DashboardView,
		meta: {
			app_view: 'menu',
			app_title: 'Dashboard',
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

router.beforeEach((to) => {
	if (typeof to.meta.app_title === 'string') {
		appSetTitle(to.meta.app_title);
	} else {
		appSetTitle(null);
	}
});

createApp(App).use(router).mount('#app-root')

enablePwaInstall();
