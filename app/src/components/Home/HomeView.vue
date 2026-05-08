<script setup lang="ts">
import { useRouter } from 'vue-router';
import GenericButton from '../App/GenericButton.vue';
import { intl, useLanguage } from '../../intl';

const router = useRouter();
const lang = useLanguage();

const openApp = () => {
	router.push('/app');
};

</script>

<template>

	<div class="home-view">

		<main>

			<header>

				<h1>
					<template v-if="lang === 'uk'">
						Нарешті, щось схоже на справжні флеш-картки
					</template>
					<template v-else-if="lang === 'de'">
						Gamifiziere das Lernerlebnis
						<br>
						mit wie-UNO Lernkarten
					</template>
					<template v-else>
						Gamify learning experience
						<br>
						with UNO-like flashcards
					</template>
				</h1>

			</header>

			<section>
				<img src="/src/assets/images/card-demo.svg" width="619" height="804" alt="Card demo image" />
			</section>

			<section>
				<GenericButton @click="openApp">
					{{ intl(lang, {
						en: 'Experience the App',
						de: 'Die App Erfahren',
						uk: 'До застосунку'
					}) }}
				</GenericButton>
			</section>

			<section>
				<p>
					{{ intl(lang, {
						en: 'Completely free of charge, no ads, self hostable by design.',
						de: 'Vollständig kostenlos, werbefrei und von Haus aus selbst hostbar.',
						uk: 'Абсолютно безкоштовно, без реклами, з можливістю самостійного розміщення.'
					}) }}
				</p>
			</section>

		</main>

		<footer>
			<p>
				Designed by maddsua. Provided by MWS via Railway Corp.
			</p>
		</footer>

	</div>

</template>

<style lang="scss" scoped>
	.home-view {
		display: flex;
		flex-direction: column;
		align-items: center;
		height: 100vh;
		height: 100dvh;
		overflow: hidden auto;
		scrollbar-width: thin;
		padding: 1rem;

		main {
			display: flex;
			flex-direction: column;
			align-items: center;
			gap: 2rem;
			padding: 2rem;
			width: 100%;
			flex-grow: 1;

			@media (orientation: landscape) {
				max-width: 50rem;
			}
		}

		h1 {
			text-align: center;
			font-size: 3rem;

			@media (orientation: portrait) {
				font-size: 2rem;
			}
		}

		p {
			text-align: center;
		}

		img {
			display: block;
			width: 20rem;
			max-width: 100%;
			height: auto;
		}

		footer {
			font-size: 0.75rem;
			color: var(--app-theme-mysterious-white);
		}
	}
</style>
