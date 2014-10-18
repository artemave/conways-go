var gulp           = require('gulp');
var pogo           = require('gulp-pogo');
var browserify     = require('browserify');
var sass           = require('gulp-sass');
var concat         = require('gulp-concat');
var plumber        = require('gulp-plumber');
var gutil          = require('gulp-util');
var fs             = require('fs');
var watch          = require('gulp-watch');
var karma          = require('karma').server;
var watchify       = require('watchify');
var source         = require('vinyl-source-stream');
var gulpBowerFiles = require('main-bower-files');

var onError = function (err) {
  gutil.beep();
  gutil.log(gutil.colors.red(err.message))
  gutil.log(err)
};

gulp.task("bower-files", function() {
  return gulp.src(gulpBowerFiles())
    .pipe(concat('deps.js'))
    .pipe(gulp.dest("./public"))
});

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

gulp.task('compile-pogo', function(callback){
    return gulp.src('./public/{js,test}/**/*.pogo')
      .pipe(plumber({errorHandler: onError}))
      .pipe(pogo())
      .pipe(gulp.dest('./public/'));
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

gulp.task("watch", ["compile-pogo", "watchify"], function() {
  watch('./public/{js,test}/**/*.pogo')
    .pipe(plumber({errorHandler: onError}))
    .pipe(pogo())
    .pipe(gulp.dest('./public/'));

  gulp.watch('./public/css/**', ['styles']);
})

gulp.task('default', ['styles', 'bower-files', 'browserify']);
