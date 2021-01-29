using Dictionary = System.Collections.Generic.Dictionary<
    DF.Game.API.Constants.EntityProperty,
    DF.Game.Curve.CurveOneOf>;

namespace DF.Game.Curve
{
    public class List
    {
        private Dictionary _curves;

        public List()
        {
            _curves = new Dictionary();
        }

        public void Add(DF.Game.Curve.CurveOneOf c)
        {
            if (_curves.ContainsKey(c.Property))
            {
                throw new System.ArgumentException(
                    System.String.Format("Cannot add {0} curve, already exists in list.", c.Property));
            }
            _curves[c.Property] = c;
        }

        public DF.Game.Curve.CurveOneOf Curve(DF.Game.API.Constants.EntityProperty p)
        {
            return _curves[p];
        }

        public System.Collections.ObjectModel.ReadOnlyCollection<DF.Game.API.Constants.EntityProperty> Properties
        {
            get => new System.Collections.ObjectModel.ReadOnlyCollection<DF.Game.API.Constants.EntityProperty>(
                new System.Collections.Generic.List<DF.Game.API.Constants.EntityProperty>(_curves.Keys)
            );
        }

        public void Merge(List cs)
        {
            foreach (var p in cs.Properties)
            {
                Curve(p).Merge(cs.Curve(p));
            }
        }
    }
}