namespace DF {
  struct datum {
    private double _tick;
    private DF.Game.API.Data.Position _value;

    public datum(double tick) {
      _tick = tick;
      _value = null;
    }

    public double Tick { get => _tick; }
    public DF.Game.API.Data.Position Value { get => _value; }

    public static bool operator <(datum a, datum b) => a.Tick < b.Tick;
    public static bool operator >(datum a, datum b) => a.Tick > b.Tick;
  }

  namespace Curve {
    class Curve {
	public static OneOf.OneOf<LinearMove> Import(DF.Game.API.Data.Curve pb) {
          switch (pb.Type) {
            case DF.Game.API.Constants.CurveType.LinearMove:
              return new LinearMove(pb.CurveId, pb.EntityId);
            default:
              break;
          }
          throw new System.ArgumentException(
            System.String.Format("Input CurveType {0} is not recognized", pb.Type));
        }
    }

    // TODO(minkezhang): Explore if we can reuse existing Go implementation via
    // the c-shared build option.
    //
    // See https://www.mono-project.com/docs/advanced/pinvoke/,
    // https://github.com/bazelbuild/rules_go/issues/54,
    // https://medium.com/learning-the-go-programming-language/calling-go-functions-from-other-languages-4c7d8bcc69bf.
    class LinearMove {
      private string _id;
      private string _entityID;
      private System.Collections.Generic.List<datum> _data;

      public LinearMove(string id, string entityID) {
        _id = id;
        _entityID = entityID;
      }

      public string ID { get => _id; }
      public string EntityID { get => _entityID; }

      public void ReplaceTail(LinearMove curve) {
        if (curve._data.Count == 0) {
          return;
        }

        var i = _data.BinarySearch(curve._data[0]);
        if (i == ~_data.Count) {
        } else {
          if (i < 0) {
            i = ~i;
          }
          _data.RemoveRange(i, _data.Count - i);
        }

        _data.AddRange(curve._data);
      }

      public DF.Game.API.Data.Position Get(double tick) {
        if (_data.Count == 0) {
          return new DF.Game.API.Data.Position{};
        }

        if (tick < _data[0].Tick) {
          return _data[0].Value;
        }
        if (_data[_data.Count - 1].Tick < tick) {
          return _data[_data.Count - 1].Value;
        }

        var i = _data.BinarySearch(new datum(tick));
        if (i == ~_data.Count) {
          return _data[~i].Value;
        }
        if (i == 0) {
           return _data[0].Value;
        }
        if (i < 0) { i = ~i; }

        var tickDelta = tick - _data[i - 1].Tick;
        return new DF.Game.API.Data.Position{
          X = _data[i - 1].Value.X + (
            _data[i].Value.X - _data[i - 1].Value.X) / (
            _data[i].Tick - _data[i - 1].Tick) * tickDelta,
          Y = _data[i - 1].Value.Y + (
            _data[i].Value.Y - _data[i - 1].Value.Y) / (
            _data[i].Tick - _data[i - 1].Tick) * tickDelta,
        };
      }
    }
  }
}
