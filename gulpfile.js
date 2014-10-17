var gulp       = require('gulp');
var pogo       = require('gulp-pogo');
var browserify = require('browserify');
var sass       = require('gulp-sass');
var concat     = require('gulp-concat');
var plumber    = require('gulp-plumber');
var gutil      = require('gulp-util');
var fs         = require('fs');
var watch      = require('gulp-watch');
var karma      = require('karma').server;
var watchify   = require('watchify');
var source     = require('vinyl-source-stream');

var onError = function (err) {
  gutil.beep();
  gutil.log(gutil.colors.red(err.message))
  gutil.log(err)
};

gulp.task('styles', function (callback) {
  return gulp.src('./public/css/app.scss')
    .pipe(plumber({
      errorHandler: onError
    }))
    .pipe(sass())
    .pipe(concat('bundle.css'))
    .pipe(gulp.dest('./public'))
});

/**
 * Run test once and exit
 */
gulp.task('test', function (done) {
  karma.start({
    configFile: __dirname + '/karma.conf.js',
    singleRun: true
  }, done);
});

/**
 * Watch for file changes and re-run tests on each change
 */
gulp.task('tdd', function (done) {
  karma.start({
    configFile: __dirname + '/karma.conf.js'
  }, done);
});

gulp.task("watchify", function() {
    browserifyAndMaybeWatchify(true)
})

gulp.task("browserify", function() {
    browserifyAndMaybeWatchify(false)
})


function browserifyAndMaybeWatchify(watch) {
  var bundler = browserify("./public/js/app.js", watchify.args)

  var bundle = function() {
    return bundler
    .bundle()
    .on('error', onError)
    .pipe(source('bundle.js'))
    .pipe(gulp.dest('./public/'));
  };

  if (watch) {
    bundler = watchify(bundler);
    bundler.on("update", bundle);
  }

  bundle()
}

gulp.task("watch", ["watchify"], function() {
  watch('./public/test/**/*.pogo')
    .pipe(plumber({errorHandler: onError}))
    .pipe(pogo())
    .pipe(gulp.dest('./public/test/'));

  watch('./public/js/**/*.pogo')
    .pipe(plumber({errorHandler: onError}))
    .pipe(pogo())
    .pipe(gulp.dest('./public/js/'));

  gulp.watch('./public/css/**', ['styles']);
})

gulp.task('default', ['styles', 'browserify']);
