var gulp        = require('gulp')
var browserify  = require('gulp-browserify');
var concat      = require('gulp-concat');
var plumber     = require('gulp-plumber')
var gutil       = require('gulp-util')
var fs          = require('fs')

var onError = function (err) {
  gutil.beep();
  gutil.log(gutil.colors.red(err.message))
  gutil.log(err)
};

gulp.task('scripts', function() {
  return gulp.src('./public/js/app.js', {read: false})
    .pipe(plumber({
      errorHandler: onError
    }))
    .pipe(browserify({
      insertGlobals: true
    }))
    .pipe(concat('bundle.js'))
    .pipe(gulp.dest('./public'))
});

gulp.watch('./public/js/**', ['scripts']);

gulp.task('default', ['scripts']);
