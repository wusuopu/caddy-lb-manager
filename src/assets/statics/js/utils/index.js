export default {
  datetimeFormat (value, format) {
    if (!value) { return '' }
    return moment(value).format(format || 'YYYY-MM-DD HH:mm:ss')
  },
}