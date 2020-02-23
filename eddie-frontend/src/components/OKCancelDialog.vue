<template>
  <v-dialog v-model="modal">
    <v-card>
      <v-card-text>
        <h2 class="title">{{title}}</h2>
      </v-card-text>
      <v-card-text class="subheading grey--text">
        {{message}}
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn flat v-on:click.native="cancel">Cancel</v-btn>
        <v-btn flat v-on:click.native="ok" class="warning">OK</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default{
  name: 'OKCancel',
  props: ['modal', 'title', 'message'],
  methods: {
    cancel () {
      this.$emit('cancel')
    },
    ok () {
      this.$emit('ok')
    }
  }
}
</script>