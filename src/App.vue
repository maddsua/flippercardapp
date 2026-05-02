<script setup lang="ts">
import { ref } from 'vue';
import CardView from './components/Cards/CardView.vue';

import type { CardNode } from './components/Cards/content';
import EndscreenView from './components/Endscreen/EndscreenView.vue';

const cards: CardNode[] = [
	{
		id: 'geo-1',
		front: {
			content: [
				{
					type: 'title',
					content: `Geography 101`
				},
				{
					type: 'textbox',
					content: [
						{
							type: 'text',
							content: `What's the caital of Germany?`
						}
					]
				},
				{
					type: 'poll',
					is_quiz: true,
					content: [
						{
							value: 'Cologne',
						},
						{
							value: 'Berlin',
							is_answer: true
						},
						{
							value: 'Bonn',
						},
					]
				}
			]
		},
		back: {
			content: [
				{
					type: 'textbox',
					content: [
						{
							type: 'text',
							content: `Berlin has been the capital of Germany again since October 3, 1990`
						}
					]
				}
			]
		}
	},
	{
		id: 'tur-1',
		front: {
			content: [
				{
					type: 'title',
					content: `Turismus 101`
				},
				{
					type: 'textbox',
					content: [
						{
							type: 'text',
							content: `In which country it is considered a grave offence to order pizza with pineapple?`
						}
					]
				},
				{
					type: 'poll',
					is_quiz: true,
					content: [
						{
							value: 'Spain',
						},
						{
							value: 'Italy',
							is_answer: true
						},
						{
							value: 'USA',
						},
					]
				}
			]
		},
		back: {
			content: [
				{
					type: 'textbox',
					content: [
						{
							type: 'text',
							content: `Itallians will kick the shit out of you if you try to deface their national dish by putting pineapple on it`
						}
					]
				}
			]
		}
	}
];

//	todo: setup proper endscreens

const isEndscreen = ref(false);

</script>

<template>
	<CardView v-if="!isEndscreen" :entries="cards" @end="isEndscreen = true" />
	<EndscreenView collectionName="Test collection" :size="2" :correct="1" :time="30" v-else />
</template>
