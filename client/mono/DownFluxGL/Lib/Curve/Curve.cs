namespace DF {
  namespace Curve {
    // Public declaration of the types of curves client currently supports.
    // Do not change the order here -- OneOf.Switch and Match, as well as
    // AsTN methods, depend on this.
    public abstract class Curve : OneOf.OneOfBase<LinearMove> {
	public static Curve Import(DF.Game.API.Data.Curve pb) {
          switch (pb.Type) {
            case DF.Game.API.Constants.CurveType.LinearMove:
              return new LinearMove(pb.EntityId, pb.Tick, pb.Data);
            default:
              break;
          }
          throw new System.ArgumentException(
            System.String.Format("Input CurveType {0} is not recognized", pb.Type));
        }
    }

    class datum : System.IComparable {
      private double _tick;
      private DF.Game.API.Data.Position _value;

      public datum(double tick) {
        _tick = tick;
        _value = null;
      }
      public datum(double tick, DF.Game.API.Data.Position value) {
        _tick = tick;
        _value = value;
      }

      public double Tick { get => _tick; }
      public DF.Game.API.Data.Position Value { get => _value; }

      public int CompareTo(object obj) {
        var d = obj as datum;
        return _tick.CompareTo(d._tick);
      }

      public override string ToString() { return _tick.ToString() + ": " + _value.ToString(); }
      public static bool operator <(datum a, datum b) => a.Tick < b.Tick;
      public static bool operator >(datum a, datum b) => a.Tick > b.Tick;
    }

    // TODO(minkezhang): Explore if we can reuse existing Go implementation via
    // the c-shared build option.
    //
    // See https://www.mono-project.com/docs/advanced/pinvoke/,
    // https://github.com/bazelbuild/rules_go/issues/54,
    // https://medium.com/learning-the-go-programming-language/calling-go-functions-from-other-languages-4c7d8bcc69bf.
    public class LinearMove : Curve {
      private string _entityID;
      private System.Collections.Generic.List<datum> _data;
      private double _tick; // Act as staleness indicator.

      public LinearMove(
        string entityID,
        double tick,
        Google.Protobuf.Collections.RepeatedField<DF.Game.API.Data.CurveDatum> data) {
        _entityID = entityID;
        _tick = tick;

        _data = new System.Collections.Generic.List<datum>();
        // Assuming data is already sorted.
        foreach (var d in data) {
          _data.Add(new datum(d.Tick, d.PositionDatum));
        }
      }

      public string EntityID { get => _entityID; }
      public double Tick { get => _tick; }

      public void ReplaceTail(LinearMove curve) {
        System.Console.Error.WriteLine("--------------------- DEBUG: REPLACETAIL _tick");
        System.Console.Error.WriteLine(_tick);
        System.Console.Error.WriteLine("--------------------- DEBUG: REPLACETAIL curve.Tick");
        System.Console.Error.WriteLine(curve.Tick);
        if (_tick > curve.Tick) {
          return;
        }

        _tick = curve.Tick;

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

        System.Console.Error.WriteLine(string.Join(", ", _data));
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
