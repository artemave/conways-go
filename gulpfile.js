var gulp         = require('gulp');
var pogo         = require('gulp-pogo');
var pogoify      = require('pogoify');
var browserify   = require('browserify');
var sass         = require('gulp-sass');
var concat       = require('gulp-concat');
var plumber      = require('gulp-plumber');
var gutil        = require('gulp-util');
var fs           = require('fs');
var watch        = require('gulp-watch');
var karma        = require('karma').server;
var watchify     = require('watchify');
var source       = require('vinyl-source-stream');
var markdownify  = require('markdownify');
var argv         = require('yargs').argv;
var uglify       = require('gulp-uglify');
var gulpif       = require('gulp-if');
var buffer       = require('vinyl-buffer');
var autoprefixer = require('gulp-autoprefixer');

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
    .pipe(autoprefixer())
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
  args = watchify.args;
  args.extensions = ['.pogo', '.md'];

  var bundler = browserify("./public/js/app.pogo", args);

  bundler.transform(pogoify);
  bundler.transform(markdownify);

  var bundle = function() {
    return bundler
      .bundle()
      .on('error', onError)
      .pipe(source('bundle.js'))
      .pipe(buffer())
      .pipe(gulpif(argv.production, uglify()))
      .pipe(gulp.dest('./public/'));
  };

  if (watch) {
    bundler = watchify(bundler);
    bundler.on("update", bundle);
    bundler.on("log", gutil.log);
  }

  bundle()
}

gulp.task("watch", ["watchify"], function() {
  gulp.watch('./public/css/**', ['styles']);
})

gulp.task('default', ['styles', 'browserify']);
