<template>
  <div></div>
</template>

<script>
import Axios from 'axios';
import useLocalStore from "@/stores/local";

export default {
  name: 'PlayerView',
  props: {
    gameCode: {
      type: Number,
      required: true,
    },
    game: {
      type: Object,
      required: false,
      default() {
        return {
          Key: '',
          SuggestionsEnabled: false,
        };
      },
    },
    player: {
      type: Object,
      required: true,
    },
    volume: {
      type: Number,
      required: false,
      default() {
        return 50;
      },
    },
  },
  data() {
    return {
      localStore: useLocalStore(),
      sockMessages: [
        { m: 'NEW_PLAYER', f: this.refreshGame },
        { m: 'PLAYER_QUIT', f: this.refreshGame },
        { m: 'PLAYER_UPDATE', f: this.playerUpdate },
        { m: 'OPTION_ADDED', f: this.refreshGame },
        { m: 'OPTION_ASSIGNED', f: this.refreshGame },
        {
          m: 'OPTION_CHANGED',
          f: (data) => {
            this.optionChanged(data.Index, data.Option);
          },
        },
        {
          m: 'PLAYER_NAME_CHANGED',
          f: (data) => {
            this.playerNameChanged(data.Index, data.Name);
          },
        },
        {
          m: 'CHARACTER_EQUIPMENT_ADDED',
          f: (data) => {
            this.characterEquipmentAdded(
              data.PlayerIndex,
              data.Equipment,
              data.UUID,
            );
          },
        },
        {
          m: 'CHARACTER_EQUIPMENT_CHANGED',
          f: (data) => {
            this.characterEquipmentChanged(
              data.PlayerIndex,
              data.EquipmentIndex,
              data.Equipment,
              data.UUID,
            );
          },
        },
        {
          m: 'CHARACTER_EQUIPMENT_REMOVED',
          f: (data) => {
            this.characterEquipmentRemoved(
              data.PlayerIndex,
              data.EquipmentIndex,
              data.UUID,
            );
          },
        },
        {
          m: 'CHARACTER_NAME_CHANGED',
          f: (data) => {
            this.characterNameChanged(data.Index, data.CharacterName, data.UUID);
          },
        },
        {
          m: 'CHARACTER_RACE_CHANGED',
          f: (data) => {
            this.characterRaceChanged(data.Index, data.CharacterRace, data.UUID);
          },
        },
        {
          m: 'EXPERIENCE_UPDATED',
          f: (data) => {
            this.experienceUpdated(data.Experience);
          },
        },
        {
          m: 'HEALTH_UPDATED',
          f: (data) => {
            this.healthUpdated(data.Health);
          },
        },
        {
          m: 'PLAYER_EXPERIENCE_CHANGED',
          f: (data) => {
            this.playerExperienceChanged(data.Index, data.Experience);
          },
        },
        {
          m: 'PLAYER_HEALTH_CHANGED',
          f: (data) => {
            this.playerHealthChanged(data.Index, data.Health);
          },
        },
        {
          m: 'QUEST_CHANGED',
          f: (data) => {
            this.questChanged(data.Quest);
          },
        },
        { m: 'OPTION_VOTE', f: this.addVote },
        { m: 'ROLL_DICE', f: this.rollDice },
        {
          m: 'SET_VOLUME',
          f: (data) => {
            this.setVolume(data.Volume);
          },
        },
        {
          m: 'NEW_SUGGESTION',
          f: (data) => {
            this.newSuggestion(data.Suggestion);
          },
        },
        {
          m: 'SUGGESTIONS_ENABLED',
          f: (data) => {
            this.setSuggestions(data.Enabled);
          },
        },
      ],
    };
  },
  methods: {
    getToken() {
      return this.localStore.managedTokens.find((token) => token.Service === 'Trudge');
    },
    refreshGame() {
      const token = this.getToken();
      if (token) {
        this.$emit('gameReloading');
        Axios({
          method: 'get',
          url: `${this.apiUri}/game`,
          headers: {
            Authorization: `Bearer ${token.AccessToken}`,
          },
        }).then((r) => {
          if (r.data.Response.Game) {
            this.$emit('gameUpdated', r.data.Response.Game);
          }
        });
      }
    },
    playerUpdate(data) {
      this.$emit('playerUpdated', data);
    },
    playerNameChanged(index, value) {
      this.$emit('playerNameChanged', index, value);
    },
    playerExperienceChanged(index, value) {
      this.$emit('playerExperienceChanged', index, value);
    },
    playerHealthChanged(index, value) {
      this.$emit('playerHealthChanged', index, value);
    },
    characterEquipmentAdded(playerIndex, value, uuid) {
      this.$emit('characterEquipmentAdded', playerIndex, value, uuid);
    },
    characterEquipmentChanged(playerIndex, equipmentIndex, value, uuid) {
      this.$emit('characterEquipmentChanged', playerIndex, equipmentIndex, value, uuid);
    },
    characterEquipmentRemoved(playerIndex, equipmentIndex, uuid) {
      this.$emit('characterEquipmentRemoved', playerIndex, equipmentIndex, uuid);
    },
    characterNameChanged(index, value, uuid) {
      this.$emit('characterNameChanged', index, value, uuid);
    },
    characterRaceChanged(index, value, uuid) {
      this.$emit('characterRaceChanged', index, value, uuid);
    },
    optionChanged(index, value) {
      this.$emit('optionChanged', index, value);
    },
    questChanged(value) {
      this.$emit('questChanged', value);
    },
    addVote(data) {
      this.emitter.emit('voteAdded', data);
    },
    experienceUpdated(value) {
      this.$emit('experienceUpdated', value);
    },
    healthUpdated(value) {
      this.$emit('healthUpdated', value);
    },
    newSuggestion(suggestion) {
      this.emitter.emit('newSuggestion', suggestion);
    },
    rollDice() {
      this.emitter.emit('rollDice');
    },
    setVolume(value) {
      this.$emit('setVolume', value);
    },
    setSuggestions(value) {
      this.$emit('setSuggestions', value);
    },
    sockMsg(msg) {
      try {
        const data = JSON.parse(msg);
        const instruction = this.sockMessages.find((m) => m.m === data.Command);
        if (instruction) {
          if (data.Data) {
            instruction.f.call(this, data.Data);
          } else {
            instruction.f.call(this);
          }
        }
      } catch (e) {
        // console.log(e);
      }
    },
  },
};

</script>

<style scoped>

</style>
