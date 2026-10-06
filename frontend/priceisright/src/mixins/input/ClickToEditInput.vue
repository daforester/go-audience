<template>
  <div @focusin="activate" @click="activate" class="clickToEditInput">
    <label v-if="label" :for="`${uid}_input`" :class="labelClass">
      {{ label }}
    </label>
    <div :class="inputClass">
      <input :id="`${uid}_input`"
             ref="input"
             v-model="inputValue"
             :aria-label="ariaLabel"
             :tabindex="tabindex"
             :type="type"
             :name="name"
             @change="(e) => { $emit('change', e) }"
             @keypress.13="blur"
             @input="$emit('update:value', updatedValue($event.target.value))"
             @focusin="showInput = true"
             @focusout="showInput = false"
      />
      <span :class="inputDisplayClass"
            aria-hidden="true"
            @focusin="activate"
            @click="activate"
      >
        {{ displayValue }}
      </span>
    </div>
  </div>
</template>

<script>

import { v4 as uid } from "uuid";

export default {
  name: 'ClickToEditInput',
  props: {
    ariaLabel: {
      type: String,
      required: false,
      default: '',
    },
    inputClass: {
      type: String,
      required: false,
      default: '',
    },
    label: {
      type: String,
      required: false,
      default: '',
    },
    labelClass: {
      type: String,
      required: false,
      default: '',
    },
    locale: {
      type: String,
      required: false,
      default: 'en-GB',
    },
    name: {
      type: String,
      required: false,
      default: '',
    },
    placeholder: {
      type: String,
      required: false,
      default: '',
    },
    tabindex: {
      type: Number,
      required: false,
      default: 0,
    },
    type: {
      type: String,
      required: false,
      default: 'text',
    },
    value: {
      required: true,
    },
  },
  data() {
    return {
      inputValue: '',
      realType: '',
      showInput: false,
      uid: 'vueInstance_',
    };
  },
  computed: {
    displayValue() {
      if (!this.value && this.placeholder) {
        return this.placeholder;
      }

      return this.value;
    },
    inputDisplayClass() {
      if (!this.value && this.placeholder) {
        return 'inputDisplay usingPlaceholder';
      }

      return 'inputDisplay';
    },
  },
  mounted() {
    this.realType = typeof this.value;
    if (typeof this.value === 'string') {
      this.inputValue = this.value;
    } else {
      this.inputValue = this.value.toString(10);
    }
    this.uid = `vueInstance_${uid()}`;
  },
  methods: {
    activate() {
      this.inputValue = this.value;
      if (this.inputValue === 0) {
        this.inputValue = '';
      }
      this.$refs.input.focus();
    },
    blur() {
      this.$refs.input.blur();
      this.$emit('blur');
    },
    focus() {
      this.$refs.input.focus();
    },
    updatedValue(value) {
      if (this.realType === 'number') {
        return parseFloat(value);
      }
      return value;
    },
  },
};

</script>

<style lang="scss">
  .clickToEditInput {
    .usingPlaceholder {
      color: lightgray;
    }
    input:focus + .inputDisplay {
      display: none;
    }
    input {
      border: 0;
      outline: none;
      position: absolute;
      left: -10000px;
      top: auto;
      width: 1px;
      height: 1px;
      overflow: hidden;
    }
    input:focus {
      position:static;
      width: 100%;
      height:auto;
    }
  }
</style>
